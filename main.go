package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define a command-line flag
	dice := flag.String("d", "d6", "Specify the type of dice to roll. Format is DX where X is an integer. Default is D6.")
	flag.Parse()
	fmt.Printf("You rolled a %s\n", *dice)
}
