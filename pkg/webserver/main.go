package webserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("TEST: " + string(body))
		lightcone := auxLightCone{}
		if err := json.Unmarshal(body, &lightcone); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, buffData := range lightcone.Buffs {
			log.Println("Buff: " + fmt.Sprintf("%+v", buffData))
		}
		log.Printf("Lightcone received: %+v", lightcone)
	})
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, _ := json.Marshal(hsrtct.LightCone{
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
	err := templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"lightcone":  string(jsonData),
		"attackTags": hsrtct.AttackTagKeys(),
		"elements":   hsrtct.ElementKeys(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lightConeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		lcs, err := db.GetLightCones()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = templates.ExecuteTemplate(w, "lightcone_list.html", lcs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("TEST: " + string(body))
		lightcone := hsrtct.LightCone{}
		if err := json.Unmarshal(body, &lightcone); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, buffData := range lightcone.Buffs {
			log.Println("Buff: " + fmt.Sprintf("%+v", buffData))
		}
		log.Printf("Lightcone received: %+v", lightcone)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func Start(port int, injectedDb database) error {
	db = injectedDb
	go loadTemplates()
	mux := setupHandlers()
	log.Printf("Server starting on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
