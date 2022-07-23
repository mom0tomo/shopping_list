package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"

	"github.com/mom0tomo/shopping-list/controllers"
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
	controllers := controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/things", controllers.GetThings(db)).Methods("GET")
	router.HandleFunc("/things/{id}", controllers.GetThing(db)).Methods("GET")
	router.HandleFunc("/things", controllers.AddThing(db)).Methods("POST")
	router.HandleFunc("/things", controllers.UpdateThing(db)).Methods("PUT")
	router.HandleFunc("/things/{id}", controllers.DeleteThing(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
