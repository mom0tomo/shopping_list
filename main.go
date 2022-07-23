package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"

	"github.com/mom0tomo/shopping-list/driver"
	"github.com/mom0tomo/shopping-list/models"
)

var things []models.Thing
var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	db = driver.ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/things", getThings(db)).Methods("GET")
	router.HandleFunc("/things/{id}", getThing(db)).Methods("GET")
	router.HandleFunc("/things", addThing(db)).Methods("POST")
	router.HandleFunc("/things", updateThing(db)).Methods("PUT")
	router.HandleFunc("/things/{id}", deleteThing(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getThings(w http.ResponseWriter, r *http.Request) {
	var thing models.Thing

	rows, err := db.Query("SELECT * FROM things")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&thing.ID, &thing.Name, &thing.Maker)
		if err != nil {
			log.Fatal(err)
		}
		things = append(things, thing)
	}
}

func getThing(w http.ResponseWriter, r *http.Request) {
	var thing models.Thing

	params := mux.Vars(r)
	id := params["id"]

	rows := db.QueryRow("SELECT * FROM things WHERE id = $1", id)
	rows.Scan(&thing.ID, &thing.Name, &thing.Maker)

}

func addThing(w http.ResponseWriter, r *http.Request) {
	var thing models.Thing
	var ThingID int

	json.NewDecoder(r.Body).Decode(&thing)

	err := db.QueryRow("INSERT INTO things(name, maker) VALUES($1, $2) RETURNING id", thing.Name, thing.Maker).Scan(&ThingID)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(ThingID)
}

func updateThing(w http.ResponseWriter, r *http.Request) {
	var thing models.Thing
	json.NewDecoder(r.Body).Decode(&thing)

	result, err := db.Exec("UPDATE things SET name=$1, maker=$2 WHERE id=$3", &thing.Name, &thing.Maker, &thing.ID)
	if err != nil {
		log.Fatal(err)
	}
	roswUpdated, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(roswUpdated)
}

func deleteThing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result, err := db.Exec("DELETE FROM things WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}
	rowsDelete, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(rowsDelete)
}
