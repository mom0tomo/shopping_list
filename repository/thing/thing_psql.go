package thingRepository

import (
	"database/sql"

	"github.com/mom0tomo/shopping-list/models"
)

type ThingRepository struct{}

func (t ThingRepository) GetThings(db *sql.DB, thing models.Thing, things []models.Thing) ([]models.Thing, error) {
	rows, err := db.Query("SELECT * FROM things")
	if err != nil {
		return []models.Thing{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&thing.ID, &thing.Name, &thing.Maker)
		if err != nil {
			return []models.Thing{}, err
		}
		things = append(things, thing)
	}
	return things, err
}

func (t ThingRepository) GetThing(db *sql.DB, thing models.Thing, id int) (models.Thing, error) {
	err := db.QueryRow("SELECT * FROM things WHERE id=$1", id).Scan(&thing.ID, &thing.Name, &thing.Maker)
	if err != nil {
		return models.Thing{}, err
	}
	return thing, err
}

func (t ThingRepository) AddThing(db *sql.DB, thing models.Thing) (int, error) {
	err := db.QueryRow("INSERT INTO things (name, maker) VALUES ($1, $2)", thing.Name, thing.Maker).Scan(&thing.ID)
	if err != nil {
		return 0, err
	}
	return thing.ID, err
}

func (t ThingRepository) UpdateThing(db *sql.DB, thing models.Thing) (int64, error) {
	result, err := db.Exec("UPDATE things SET name=$1, maker=$2 WHERE id=$3", thing.Name, thing.Maker, thing.ID)
	if err != nil {
		return 0, err
	}
	roswUpdated, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return roswUpdated, nil
}

func (t ThingRepository) DeleteThing(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM things WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsDeleted, nil
}
