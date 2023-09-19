package databasesetup

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func DatabaseSetup() (*sql.DB, error) {
	conninfo := "user=postgres password='' dbname=EnterpriseNotes host=127.0.0.1 port=5432 sslmode=disable"

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.Println(err)
	}

	dbName := "EnterpriseNotes"
	_, err = db.Exec("create database " + dbName)
	if err != nil {
		//handle the error
		log.Println(err)
	}

	fmt.Println("Database Setup Complete")
	return db, nil
}
