package driver

import (
	"database/sql"
	"log"
	"os"

	pg "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	pgUrl, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", pgUrl)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
