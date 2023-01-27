package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Welcome to Juncoin CLI\n")
	fmt.Printf("Instructions:\n")
	fmt.Printf("html: 	Start HTML explorer\n")
	fmt.Printf("rest:	Start REST APIs (recommended)\n")
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "html":
		fmt.Println("Starting HTML explorer")
	case "rest":
		fmt.Println("Starting REST APIs")
	default:
		usage()
	}
}
