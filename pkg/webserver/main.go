package webserver

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

var templates = template.Must(template.New("").Parse(""))

func loadTemplates() {
	templateDir := "web/template"
	filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			templates = template.Must(templates.ParseFiles(path))
		}
		return nil
	})
}

func setupHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", helloHandler)
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", hsrtct.LightCone{
		ID:      "UUID",
		Name:    "On the Fall of an Aeon",
		Level:   80,
		BaseHp:  1058,
		BaseAtk: 529,
		BaseDef: 396,
		Buffs: []hsrtct.Buff{
			{Stat: hsrtct.AtkPct, Value: 16 * 4},
			{Stat: hsrtct.DmgBonus, Value: 24},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lightConeHandler(w http.ResponseWriter, r *http.Request) {

}

func Start(port int, injectedDb database) error {
	db = injectedDb
	loadTemplates()
	mux := setupHandlers()
	log.Printf("Server starting on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
