package webserver

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var db database

func helloHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.New("index.html")
	tmpl, _ = tmpl.ParseFiles("web/templates/index.html")
	err := tmpl.Execute(w, "Hello World!")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Setup(injectedDb database) {
	db = injectedDb
	http.HandleFunc("/", helloHandler)
}

func Start(port int) error {
	log.Printf("Server starting on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
