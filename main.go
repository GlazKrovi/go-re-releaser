package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: gorr release <patch|minor|major>")
		fmt.Println("Example: gorr release patch")
		os.Exit(1)
	}

	if os.Args[1] != "release" {
		fmt.Printf("Error: Unknown command '%s'. Use 'gorr release <patch|minor|major>'\n", os.Args[1])
		os.Exit(1)
	}

	releaseType := os.Args[2]

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

	// Tag and push the new version
	tagAndPush(nextVersion)
	fmt.Printf("Next version pushed: %s\n", nextVersion)

	// If everything is ok, create and send the complete release
	releaseOutput := callReleaser()
	fmt.Println(releaseOutput)
}

func getCurrentVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last tag: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
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
	parts := strings.Split(currentVersion, ".")
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

func callReleaser() string {
	cmd := exec.Command("goreleaser", "release", "--clean")
	output, _ := cmd.Output()
	return string(output)
}
