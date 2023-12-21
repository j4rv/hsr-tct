package webserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
)

func setupHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/lightcones", lightConeHandler)
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", hsrtct.LightCone{
		ID:      0,
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
	lcs, err := db.GetLightCones()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "lightcone_list.html", lcs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Start(port int, injectedDb database) error {
	db = injectedDb
	go loadTemplates()
	mux := setupHandlers()
	log.Printf("Server starting on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
