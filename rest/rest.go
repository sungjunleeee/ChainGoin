package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sungjunleeee/ChainGoin/blockchain"
	"github.com/sungjunleeee/ChainGoin/utils"
	"github.com/sungjunleeee/ChainGoin/wallet"
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

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type walletResponse struct {
	Address string `json:"address"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
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
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "See balance of an address",
		},
	}
	json.NewEncoder(w).Encode(data)
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(blockchain.GetAllBlocks(blockchain.Blockchain()))
	case "POST":
		blockchain.Blockchain().AddBlock()
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

func getBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	if total == "true" {
		balance := blockchain.GetBalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(w).Encode(balanceResponse{address, balance})
	} else {
		err := json.NewEncoder(w).Encode(blockchain.FilterUTxOutsByAddress(address, blockchain.Blockchain()))
		utils.HandleErr(err)
	}
}

func getMempool(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(blockchain.Mempool.Txs)
	utils.HandleErr(err)
}

func createTxs(w http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	utils.HandleErr(err)
	err = blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getWallet(w http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(w).Encode(walletResponse{Address: address})
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
	router.HandleFunc("/balance/{address}", getBalance).Methods("GET")
	router.HandleFunc("/mempool", getMempool).Methods("GET")
	router.HandleFunc("/wallet", getWallet).Methods("GET")
	router.HandleFunc("/transactions", createTxs).Methods("POST")
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
