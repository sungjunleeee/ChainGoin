package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sungjunleeee/juncoin/blockchain"
)

const (
	port         string = ":4000"
	templatePath string = "templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(w, "home", data)
}

func main() {
	templates = template.Must(template.ParseGlob(templatePath + "pages/*.gohtml"))
	// updating the existing templates variable to include the partials
	templates = template.Must(templates.ParseGlob(templatePath + "partials/*.gohtml"))
	http.HandleFunc("/", handleHome)
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
