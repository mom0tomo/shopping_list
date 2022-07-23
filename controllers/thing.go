package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mom0tomo/shopping-list/models"
	thingRepository "github.com/mom0tomo/shopping-list/repository/thing"
	"github.com/mom0tomo/shopping-list/utils"
)

type Controller struct{}

var things []models.Thing

func (c Controller) GetThings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing models.Thing
		var error models.Error

		things = []models.Thing{}
		thingRepo := thingRepository.ThingRepository{}
		things, err := thingRepo.GetThings(db, thing, things)
		if err != nil {
			error.Message = "Server Error	"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, things)

		json.NewEncoder(w).Encode(thing)
	}
}

func (c Controller) GetThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing models.Thing
		var error models.Error

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		thingRepo := thingRepository.ThingRepository{}

		thing, err := thingRepo.GetThing(db, thing, id)
		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "No Thing Found"
				utils.SendError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server Error	"
				utils.SendError(w, http.StatusInternalServerError, error)
			}
			w.Header().Set("Content-Type", "application/json")
			utils.SendSuccess(w, thing)

			json.NewEncoder(w).Encode(thing)
		}
	}
}

func (c Controller) AddThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing models.Thing
		var error models.Error

		json.NewDecoder(r.Body).Decode(&thing)

		if thing.Name == "" || thing.Maker == "" {
			error.Message = "Enter missing fields"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		thingRepo := thingRepository.ThingRepository{}
		thingID, err := thingRepo.AddThing(db, thing)
		if err != nil {
			error.Message = "Server Error	"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, thingID)
	}
}

func (c Controller) UpdateThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var thing models.Thing
		var error models.Error

		json.NewDecoder(r.Body).Decode(&thing)
		if thing.ID == 0 || thing.Name == "" || thing.Maker == "" {
			error.Message = "All Fields are required"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		thingRepo := thingRepository.ThingRepository{}
		rowsUpdated, err := thingRepo.UpdateThing(db, thing)
		if err != nil {
			error.Message = "Server Error	"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (c Controller) DeleteThing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)

		thingRepo := thingRepository.ThingRepository{}
		id, _ := strconv.Atoi(params["id"])

		rowsDeleted, err := thingRepo.DeleteThing(db, id)
		if err != nil {
			error.Message = "Server Error	"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		if rowsDeleted == 0 {
			error.Message = "No Thing Found"
			utils.SendError(w, http.StatusNotFound, error)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, rowsDeleted)
		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
