package main

import (
	"EnterpriseNotes/databasesetup"
)

func main() {
	db, err := databasesetup.DatabaseSetup()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
