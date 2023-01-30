package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sungjunleeee/juncoin/blockchain"
	"github.com/sungjunleeee/juncoin/utils"
)

var port string

// URL is a custom type for overriding methods (MarshalText)
type url string

// MarshalText implements to intervenes in the json.Marshal process
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

// addBlockBody is a struct for POST /blocks
type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func showDocumentation(w http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See status of the blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See all blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a block",
		},
	}
	json.NewEncoder(w).Encode(data)
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(blockchain.Blockchain().GetAllBlocks())
	case "POST":
		var addBlockBody addBlockBody
		// Decode (Unmarshal) is case-insensitive: Message == message
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.Blockchain().AddBlock(addBlockBody.Message)
		w.WriteHeader(http.StatusCreated)
	}
}

func findBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(w)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{err.Error()})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain.Blockchain())
}

// Start starts the rest API
func Start(newPort int) {
	port = fmt.Sprintf(":%d", newPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware) // middleware for routes below
	router.HandleFunc("/", showDocumentation).Methods("GET")
	router.HandleFunc("/status", getStatus).Methods("GET")
	router.HandleFunc("/blocks", getBlocks).Methods("GET", "POST") // won't allow other methods
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", findBlock).Methods("GET")
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
