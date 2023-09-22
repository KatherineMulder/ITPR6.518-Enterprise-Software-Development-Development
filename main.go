package main

import (
	"fmt"
	"log"
	"context"

	"EnterpriseNotes/databasesetup"
	
)

func main() {
	conn, err := databasesetup.DatabaseSetup()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(context.Background()) // Close the database connection when done

	fmt.Println("Database setup completed successfully.")

	// Example of creating a new note
	err = databasesetup.CreateNote(conn, "Sample Note Title", "Sample Note Content", "Active", 1)
	if err != nil {
		log.Fatal("Failed to create a new note: ", err)
	}

	fmt.Println("Note created successfully.")

	// retrieving a note by ID
	note, err := databasesetup.GetNoteByID(conn, 1) 
	if err != nil {
		log.Fatal("Failed to retrieve a note: ", err)
	}

	fmt.Printf("Note ID: %d\nTitle: %s\nContent: %s\nStatus: %s\nUserID: %d\n",
		note.ID, note.Title, note.Content, note.Status, note.UserID)

	// updating a note
	err = databasesetup.UpdateNote(conn, 1, "Updated Title", "Updated Content", "Inactive") 
	if err != nil {
		log.Fatal("Failed to update the note: ", err)
	}

	fmt.Println("Note updated successfully.")

	// deleting a note by ID
	err = databasesetup.DeleteNoteByID(conn, 1) 
	if err != nil {
		log.Fatal("Failed to delete the note: ", err)
	}

	fmt.Println("Note deleted successfully.")
}
