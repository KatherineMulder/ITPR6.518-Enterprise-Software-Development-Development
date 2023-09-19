package main

import (
	"fmt"
	"log"

	"EnterpriseNotes/databasesetup"
)

func main() {
	_, err := databasesetup.DatabaseSetup()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database setup completed successfully.")
}