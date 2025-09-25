package main

import (
	"fmt"
	"os"
)

func main() {
	// Check if we have the right number of arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: gorr release <patch|minor|major>")
		fmt.Println("Example: gorr release patch")
		os.Exit(1)
	}

	// Check if the first argument is "release"
	if os.Args[1] != "release" {
		fmt.Printf("Error: Unknown command '%s'. Use 'gorr release <patch|minor|major>'\n", os.Args[1])
		os.Exit(1)
	}

	// Get the release type from the second argument
	releaseType := os.Args[2]

	// Validate the release type
	validTypes := map[string]bool{
		"patch": true,
		"minor": true,
		"major": true,
	}

	if !validTypes[releaseType] {
		fmt.Printf("Error: Invalid release type '%s'. Must be one of: patch, minor, major\n", releaseType)
		os.Exit(1)
	}

	// Display the chosen option
	fmt.Printf("Release type chosen: %s\n", releaseType)
}
