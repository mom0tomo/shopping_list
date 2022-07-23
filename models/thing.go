package models

import (
	"database/sql"
)

type Thing struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Maker string `json:"maker"`
}

var things []Thing
var db *sql.DB
