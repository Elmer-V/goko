package main

import (
	"fmt"
	"os"
)

const baseURL = "https://jutge.org/problems/"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: goko <mode> <id>")
		os.Exit(1)
	}

	mode := os.Args[1]
	switch mode {
	case "set", "s":
		set()
	case "check", "c":
		check()

	default:
		fmt.Println("Not a valid argument. Use set/s or check/c instead")
	}

}
