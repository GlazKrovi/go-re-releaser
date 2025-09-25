package main

import (
	"fmt"
	"os"
	"os/exec"
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
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing goreleaser: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("GoReleaser version: %s", string(output))
}
