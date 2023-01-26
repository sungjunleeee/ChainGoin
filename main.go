package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":4000"

// URL is a custom type for overriding methods (MarshalText)
type URL string

// MarshalText implements to intervenes in the json.Marshal process
func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// URLDescription is a struct for API calls
type URLDescription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "See documentation",
		},
		{
			URL:         URL("/blocks"),
			Method:      "POST",
			Description: "See documentation",
			Payload:     "data:string",
		},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	// equivalent to:
	// b, err := json.Marshal(data)
	// utils.handleErr(err)
	// fmt.Fprintf(w, "%s", b)
}

func main() {
	// The below will be called twice
	// since the browser will send a request for the favicon
	http.HandleFunc("/", documentation)
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
