package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Thing struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Maker string `json:"maker"`
}

var things []Thing

func main() {
	router := mux.NewRouter()

	things = append(things,
		Thing{ID: 1, Name: "Thing 1", Maker: "Maker 1"},
		Thing{ID: 2, Name: "Thing 2", Maker: "Maker 2"},
		Thing{ID: 3, Name: "Thing 3", Maker: "Maker 3"},
		Thing{ID: 4, Name: "Thing 4", Maker: "Maker 4"},
		Thing{ID: 5, Name: "Thing 5", Maker: "Maker 5"},
	)

	router.HandleFunc("/things", getThings).Methods("GET")
	router.HandleFunc("/things/{id}", getThing).Methods("GET")
	router.HandleFunc("/things", addThing).Methods("POST")
	router.HandleFunc("/things", updateThing).Methods("PUT")
	router.HandleFunc("/things/{id}", deleteThing).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getThings(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(things)
}

func getThing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	for _, thing := range things {
		if thing.ID == id {
			json.NewEncoder(w).Encode(&thing)
			return
		}
	}
}

func addThing(w http.ResponseWriter, r *http.Request) {
	var thing Thing
	_ = json.NewDecoder(r.Body).Decode(&thing)

	things = append(things, thing)
	json.NewEncoder(w).Encode(things)
}

func updateThing(w http.ResponseWriter, r *http.Request) {
	var thing Thing
	json.NewDecoder(r.Body).Decode(&thing)

	for i, t := range things {
		if t.ID == thing.ID {
			things[i] = thing
		}
	}
	json.NewEncoder(w).Encode(things)
}

func deleteThing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}
	for i, t := range things {
		if t.ID == id {
			things = append(things[:i], things[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(things)
}
