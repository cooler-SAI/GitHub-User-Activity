package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}

	username := os.Args[1]
	fmt.Printf("Fetching activity for user: %s\n", username)
}
