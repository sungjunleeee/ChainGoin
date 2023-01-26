package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/sungjunleeee/juncoin/blockchain"
)

const (
	port         string = ":4000"
	templatePath string = "explorer/templates/"
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

func handleAdd(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

}

// Start starts the web server
func Start() {
	templates = template.Must(template.ParseGlob(templatePath + "pages/*.gohtml"))
	// updating the existing templates variable to include the partials
	templates = template.Must(templates.ParseGlob(templatePath + "partials/*.gohtml"))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/add", handleAdd)
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
