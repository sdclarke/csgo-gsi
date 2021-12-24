package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sdclarke/csgo-gsi/pkg/structs"
)

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		d := json.NewDecoder(r.Body)
		var i structs.State
		d.Decode(&i)
		if i.Provider != nil {
			log.Printf("Provider: %#v", *i.Provider)
		}
		if i.Map != nil {
			log.Printf("Map: %#v", *i.Map)
		}
		if i.Round != nil {
			log.Printf("Round: %#v", *i.Round)
		}
		if i.Player != nil {
			log.Printf("Player: %#v", *i.Player)
		}
		if i.Previously != nil {
			log.Printf("Previously: %#v", *i.Previously)
		}
		if i.Added != nil {
			log.Printf("Added: %#v", *i.Added)
		}
		if i.AllPlayers != nil {
			log.Printf("All Players: %#v", i.AllPlayers)
		}
		if i.Auth != nil {
			log.Printf("Auth: %#v", *i.Auth)
		}
		w.WriteHeader(http.StatusOK)
	default:
		fmt.Fprintf(w, "No\n")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func main() {
	http.HandleFunc("/", handle)
	//dir := http.Dir(os.Getenv("HOME"))
	//fs := http.FileServer(dir)
	//http.Handle("/", fs)
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
	select {}
}
