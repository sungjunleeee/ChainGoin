package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sungjunleeee/juncoin/blockchain"
)

const port string = ":4000"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	data := homeData{"Home", blockchain.GetBlockChain().AllBlocks()}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handleHome)
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
