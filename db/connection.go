package db

import (
	"database/sql"
)

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/condogenius")
	if err != nil {
		return nil, err
	}
	return db, nil
}
