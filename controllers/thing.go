package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mom0tomo/shopping-list/models"
	"github.com/mom0tomo/shopping-list/utils"
)

type Controller struct{}

var things []models.Thing

func (c Controller) GetThings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing models.Thing
		var error models.Error
		json.NewDecoder(r.Body).Decode(&thing)

		things = []models.Thing{}
		rows, err := db.Query("SELECT * FROM things")
		if err != nil {
			error.Message = "Server Error	"
			utils.SendError(w, http.StatusInternalServerError, error)
		}
		defer rows.Close()
		json.NewEncoder(w).Encode(thing)
	}
}

func (c Controller) GetThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing models.Thing

		params := mux.Vars(r)
		id := params["id"]

		rows := db.QueryRow("SELECT * FROM things WHERE id=$1", id)
		rows.Scan(&thing.ID, &thing.Name, &thing.Maker)
		json.NewEncoder(w).Encode(thing)
	}
}

func (c Controller) AddThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var thing models.Thing
		var ThingID int

		json.NewDecoder(r.Body).Decode(&thing)

		err := db.QueryRow("INSERT INTO things(name, maker) VALUES($1, $2) RETURNING id", thing.Name, thing.Maker).Scan(&ThingID)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(ThingID)
	}
}

func (c Controller) UpdateThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing models.Thing

		json.NewDecoder(r.Body).Decode(&thing)

		result, err := db.Exec("UPDATE things SET name=$1, maker=$2 WHERE id=$3 RETURNING id", &thing.Name, &thing.Maker, &thing.ID)
		if err != nil {
			log.Fatal(err)
		}
		roswUpdated, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(roswUpdated)
	}
}

func (c Controller) DeleteThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
}
