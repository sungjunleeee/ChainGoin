package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/sungjunleeee/ChainGoin/blockchain"
)

const (
	templatePath string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.Blockchain().GetAllBlocks()}
	templates.ExecuteTemplate(w, "home", data)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "add", nil)
	case "POST":
		blockchain.Blockchain().AddBlock()
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

}

// Start starts the web server
func Start(newPort int) {
	templates = template.Must(template.ParseGlob(templatePath + "pages/*.gohtml"))
	// updating the existing templates variable to include the partials
	templates = template.Must(templates.ParseGlob(templatePath + "partials/*.gohtml"))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/add", handleAdd)
	fmt.Printf("Server is running on http://localhost:%d\n", newPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", newPort), nil))
}
