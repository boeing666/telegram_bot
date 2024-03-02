package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Init(query string) (*sql.DB, error) {
	db, err := sql.Open("mysql", query)
	if err != nil {
		return nil, err
	}
	return db, nil
}
