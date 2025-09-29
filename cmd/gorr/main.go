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

// Version information
const (
	Version = "1.0.1"
)

func main() {

	// Handle --version flag
	if len(os.Args) >= 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		showVersion()
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: gorr release <patch|minor|major> [args...]")
		fmt.Println("  📤 Officially release on remote repository: gorr release patch")
		fmt.Println("  🧪 Create locally: gorr release patch --snapshot")
		os.Exit(1)
	}

	if os.Args[1] != "release" {
		// Pass all arguments to goreleaser directly
		err := callGoreleaserDirect(os.Args[1:])
		if err != nil {
			fmt.Printf("❌ GoReleaser failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	releaseType := os.Args[2]

	// Capture additional arguments for goreleaser
	var goreleaserArgs []string
	if len(os.Args) > 3 {
		goreleaserArgs = os.Args[3:]
	}

	validTypes := map[string]bool{
		"patch": true,
		"minor": true,
		"major": true,
	}

	if !validTypes[releaseType] {
		fmt.Printf("Error: Invalid release type '%s'. Must be one of: patch, minor, major\n", releaseType)
		os.Exit(1)
	}

	gitErr := checkGitStatus()
	if gitErr != nil {
		fmt.Println("❌ Please commit or stash your working tree before creating a new version")
		os.Exit(1)
	}

	// Verify that goreleaser is installed
	cmd := exec.Command("goreleaser", "--version")
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing goreleaser: %v\n", err)
		os.Exit(1)
	}

	// Get the current version
	currentVersion, err := getCurrentVersion()
	if err != nil {
		fmt.Printf("Error getting current version: %v\n", err)
		os.Exit(1)
	}
	if !isValidVersionTag(currentVersion) {
		fmt.Printf("Error: Invalid version tag format '%s'. Expected format: vx.x.x (e.g., v1.2.3)\n", currentVersion)
		os.Exit(1)
	}

	// Get next version according to the release type
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

		// Tag and push the new version
		tagAndPush(nextVersion)
		fmt.Printf("Next version pushed: %s\n", nextVersion)
	}

	// If everything is ok, create and send the complete release
	err = callReleaser(goreleaserArgs)
	if err != nil {
		fmt.Printf("❌ Release failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Release completed successfully!")
}

func getCurrentVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		// If no tags exist, return default version v0.1.0
		return "v0.1.0", nil
	}

	version := strings.TrimSpace(string(output))

	// If output is empty, return default version
	if version == "" {
		return "v0.1.0", nil
	}

	// Handle snapshot versions like "v1.0.0-7-g73abd8e" -> extract "v1.0.0"
	if strings.Contains(version, "-") {
		parts := strings.Split(version, "-")
		version = parts[0]
	}

	return version, nil
}

// Check if the tag follows the vx.x.x format
func isValidVersionTag(tag string) bool {
	pattern := `^v\d+\.\d+\.\d+$`
	matched, err := regexp.MatchString(pattern, tag)
	if err != nil {
		return false
	}
	return matched
}

func getNextVersion(currentVersion string, releaseType string) string {
	// Remove 'v' prefix before parsing
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

	// Redirect stdout and stderr to the current process streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to complete
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git push failed: %v", err)
	}

	return nil
}

func tagAndPush(version string) {
	cmd := exec.Command("git", "tag", version)
	cmd.Output()

	cmd = exec.Command("git", "push", "origin", version)
	cmd.Output()
}

// Check git status before release
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
	// Build the command with base arguments
	cmdArgs := []string{"release", "--clean"}

	// Add user-provided arguments
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("goreleaser", cmdArgs...)

	// Redirect stdout and stderr to the current process streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to complete
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("goreleaser failed: %v", err)
	}

	return nil
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// callGoreleaserDirect passes all arguments directly to goreleaser
func callGoreleaserDirect(args []string) error {
	cmd := exec.Command("goreleaser", args...)

	// Redirect stdout and stderr to the current process streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to complete
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

// showVersion displays version information
func showVersion() {
	fmt.Printf("GORR - Go-RE-Releaser v%s\n", Version)
}
