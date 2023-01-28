package main

import (
	"flag"
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

	rest := flag.NewFlagSet("rest", flag.ExitOnError)
	portFlag := rest.Int("port", 4000, "Sets the port of the server")

	switch os.Args[1] {
	case "html":
		fmt.Println("Starting HTML explorer")
	case "rest":
		rest.Parse(os.Args[2:])
	default:
		usage()
	}

	if rest.Parsed() {
		fmt.Println(portFlag)
		fmt.Println("Start server")
	}
}
