package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See documentation",
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
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See a block",
		},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	// equivalent to:
	// b, err := json.Marshal(data)
	// utils.handleErr(err)
	// fmt.Fprintf(w, "%s", b)
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blockchain.GetBlockChain().AllBlocks())
	case "POST":
		var addBlockBody addBlockBody
		// Decode (Unmarshal) is case-insensitive: Message == message
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockChain().AddBlock(addBlockBody.Message)
		w.WriteHeader(http.StatusCreated)
	}
}

func findBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block, err := blockchain.GetBlockChain().FindBlock(id)
	encoder := json.NewEncoder(w)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{err.Error()})
	} else {
		encoder.Encode(block)
	}
}

// Start starts the rest API
func Start(newPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", newPort)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", getBlocks).Methods("GET", "POST") // won't allow other methods
	router.HandleFunc("/blocks/{height:[0-9]+}", findBlock).Methods("GET")
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
