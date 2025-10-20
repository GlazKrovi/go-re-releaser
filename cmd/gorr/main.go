package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  gorr release <local|patch|minor|major> [args...]")
		fmt.Println("  📤 Officially release on remote repository: gorr release patch")
		fmt.Println("  🧪 Create locally: gorr release patch --snapshot")
		fmt.Println("  🗑️  Undo last release (DEV ONLY): gorr unrelease")
		os.Exit(1)
	}

	if os.Args[1] == "unrelease" {
		err := unreleaseCommand()
		if err != nil {
			fmt.Printf("❌ Unrelease failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if os.Args[1] == "release" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: gorr release <local|patch|minor|major> [args...]")
			fmt.Println("  📤 Officially release on remote repository: gorr release patch")
			fmt.Println("  🧪 Create locally: gorr release patch --snapshot")
			os.Exit(1)
		}
	} else {
		err := callGoreleaserDirect(os.Args[1:])
		if err != nil {
			fmt.Printf("❌ GoReleaser failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	releaseType := os.Args[2]

	var goreleaserArgs []string
	if len(os.Args) > 3 {
		goreleaserArgs = os.Args[3:]
	}

	validTypes := map[string]bool{
		"local": true,
		"patch": true,
		"minor": true,
		"major": true,
	}

	if !validTypes[releaseType] {
		fmt.Printf("Error: Invalid release type '%s'. Must be one of: local, patch, minor, major\n", releaseType)
		os.Exit(1)
	}

	gitErr := checkGitStatus()
	if gitErr != nil {
		fmt.Println("❌ Please commit or stash your working tree before creating a new version")
		os.Exit(1)
	}

	if releaseType == "local" {
		fmt.Println("🧪 Creating local snapshot release...")
		err := callReleaser([]string{"--snapshot"})
		if err != nil {
			fmt.Printf("❌ Local release failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("• ...and go-re-releaser!")
		return
	}

	cmd := exec.Command("goreleaser", "--version")
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing goreleaser: %v\n", err)
		os.Exit(1)
	}

	currentVersion, err := getCurrentVersion()
	if err != nil {
		fmt.Printf("Error getting current version: %v\n", err)
		os.Exit(1)
	}
	if !isValidVersionTag(currentVersion) {
		fmt.Printf("Error: Invalid version tag format '%s'. Expected format: vx.x.x (e.g., v1.2.3)\n", currentVersion)
		os.Exit(1)
	}

	nextVersion := getNextVersion(currentVersion, releaseType)

	isSnapshot := contains(goreleaserArgs, "--snapshot")
	if isSnapshot {
		fmt.Println("Any commits would be pushed to the repo.")
		fmt.Printf("Next version that would be pushed: %s\n", nextVersion)
	} else {
		fmt.Printf("Pushing changes to remote repository...\n")
		err = gitPushChanges()
		if err != nil {
			os.Exit(1)
		}

		err = tagAndPush(nextVersion)
		if err != nil {
			fmt.Printf("❌ Failed to tag and push: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Next version pushed: %s\n", nextVersion)
	}

	err = callReleaser(goreleaserArgs)
	if err != nil {
		fmt.Printf("❌ Release failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("• ...and go-re-releaser!")
}

// unreleaseCommand removes the latest tag (local and remote) and the dist folder
func unreleaseCommand() error {
	fmt.Println("⚠️  WARNING: This command should ONLY be used in DEVELOPMENT!")
	fmt.Println("⚠️  This will delete the latest tag (local AND remote) and the dist/ folder.")
	fmt.Print("Are you sure you want to continue? (y/N): ")

	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("❌ Unrelease cancelled")
		return nil
	}

	// Get the latest local tag
	localTag, err := getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get local tag: %v", err)
	}
	if localTag == "v0.1.0" {
		return errors.New("no tags found to delete")
	}

	// Get the latest remote tag
	fetchCmd := exec.Command("git", "fetch", "--tags")
	fetchCmd.Run()

	remoteCmd := exec.Command("git", "ls-remote", "--tags", "origin")
	remoteOutput, err := remoteCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get remote tags: %v", err)
	}

	// Extract the latest remote tag
	remoteTag := ""
	lines := strings.Split(string(remoteOutput), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			ref := parts[1]
			if strings.HasPrefix(ref, "refs/tags/") {
				tag := strings.TrimPrefix(ref, "refs/tags/")
				// Skip annotated tag refs (^{})
				if !strings.HasSuffix(tag, "^{}") {
					remoteTag = tag
					break
				}
			}
		}
	}

	// Verify local and remote tags match
	if localTag != remoteTag {
		fmt.Printf("⚠️  Warning: Local tag (%s) doesn't match remote tag (%s)\n", localTag, remoteTag)
		fmt.Print("Do you want to continue anyway? (y/N): ")
		var response2 string
		fmt.Scanln(&response2)
		response2 = strings.ToLower(strings.TrimSpace(response2))
		if response2 != "y" && response2 != "yes" {
			return errors.New("tag mismatch - operation cancelled")
		}
	}

	fmt.Printf("🗑️  Deleting tag: %s\n", localTag)

	// Delete local tag
	delLocalCmd := exec.Command("git", "tag", "-d", localTag)
	delLocalCmd.Stdout = os.Stdout
	delLocalCmd.Stderr = os.Stderr
	err = delLocalCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete local tag: %v", err)
	}

	// Delete remote tag
	delRemoteCmd := exec.Command("git", "push", "origin", "--delete", localTag)
	delRemoteCmd.Stdout = os.Stdout
	delRemoteCmd.Stderr = os.Stderr
	err = delRemoteCmd.Run()
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to delete remote tag: %v\n", err)
	} else {
		fmt.Printf("✅ Remote tag deleted: %s\n", localTag)
	}

	// Delete dist folder
	fmt.Println("🗑️  Deleting dist/ folder...")
	rmCmd := exec.Command("rm", "-rf", "dist")
	err = rmCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete dist folder: %v", err)
	}

	fmt.Println("✅ Unrelease completed successfully!")
	return nil
}

func getCurrentVersion() (string, error) {
	fetchCmd := exec.Command("git", "fetch", "--tags")
	fetchCmd.Run()

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "v0.1.0", nil
	}

	version := strings.TrimSpace(string(output))

	if version == "" {
		return "v0.1.0", nil
	}

	// extract base version from snapshot versions (e.g. v1.0.0-7-g73abd8e -> v1.0.0)
	if strings.Contains(version, "-") {
		parts := strings.Split(version, "-")
		version = parts[0]
	}

	return version, nil
}

func isValidVersionTag(tag string) bool {
	pattern := `^v\d+\.\d+\.\d+$`
	matched, err := regexp.MatchString(pattern, tag)
	if err != nil {
		return false
	}
	return matched
}

func getNextVersion(currentVersion string, releaseType string) string {
	version := strings.TrimPrefix(currentVersion, "v")
	parts := strings.Split(version, ".")

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	switch releaseType {
	case "patch":
		patch++
	case "minor":
		minor++
		patch = 0
	case "major":
		major++
		minor = 0
		patch = 0
	}

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
}

func gitPushChanges() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git push failed: %v", err)
	}

	return nil
}

func tagAndPush(version string) error {
	delCmd := exec.Command("git", "tag", "-d", version)
	delCmd.Run()

	cmd := exec.Command("git", "tag", "-a", version, "-m", fmt.Sprintf("Release %s", version))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tag %s: %v", version, err)
	}

	cmd = exec.Command("git", "push", "origin", version, "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push tag %s: %v", version, err)
	}

	return nil
}

func checkGitStatus() error {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	if len(output) > 0 {
		return errors.New("git working tree is dirty")
	}
	return nil
}

func callReleaser(args []string) error {
	cmdArgs := []string{"release", "--clean"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("goreleaser", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("goreleaser failed: %v", err)
	}

	return nil
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func callGoreleaserDirect(args []string) error {
	cmd := exec.Command("goreleaser", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
