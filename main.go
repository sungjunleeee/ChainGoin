package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sungjunleeee/juncoin/utils"
)

const port string = ":4000"

type URLDescription struct {
	URL         string
	Method      string
	Description string
}

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See documentation",
		},
	}
	// convert interface to JSON
	b, err := json.Marshal(data)
	utils.HandleErr(err)
	fmt.Printf("%s\n", b)
}

func main() {
	// The below will be called twice
	// since the browser will send a request for the favicon
	http.HandleFunc("/", documentation)
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
