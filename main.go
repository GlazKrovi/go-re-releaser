package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

	fmt.Printf("Release type chosen: %s\n", releaseType)

	// Verify that goreleaser is installed
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

	fmt.Printf("Current version: %s\n", currentVersion)
}

func getCurrentVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last tag: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// isValidVersionTag checks if the tag follows the vx.x.x format
func isValidVersionTag(tag string) bool {
	// Regex pattern for vx.x.x format (e.g., v1.2.3, v0.1.0, v10.20.30)
	pattern := `^v\d+\.\d+\.\d+$`
	matched, err := regexp.MatchString(pattern, tag)
	if err != nil {
		return false
	}
	return matched
}
