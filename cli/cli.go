package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/sungjunleeee/juncoin/explorer"
	"github.com/sungjunleeee/juncoin/rest"
)

func usage() {
	fmt.Printf("Welcome to Juncoin CLI\n")
	fmt.Printf("Instructions for flags:\n")
	fmt.Printf("-mode:	Choose between html | rest | all\n")
	fmt.Printf("-port:	Set the port of the server\n")
	os.Exit(1)
}

// Start starts the cli
func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between html | rest | all")

	flag.Parse()

	switch *mode {
	// start rest api
	case "rest":
		rest.Start(*port)
	// start html explorer
	case "html":
		explorer.Start(*port)
	case "all":
		go rest.Start(*port)
		explorer.Start(*port + 1)
	default:
		usage()
	}
}
