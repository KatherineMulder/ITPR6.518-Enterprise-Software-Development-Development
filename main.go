package main

func main() {
	a := App{}
	a.Initialize()
	a.Run("")
}

/*
import (
	"context"
	"fmt"
	"log"
	"time"

	"EnterpriseNotes/databasesetup"
	"EnterpriseNotes/managenote"
	"EnterpriseNotes/models"
)

func main() {
	conn, err := databasesetup.DatabaseSetup()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(context.Background()) // Close the database connection when done

	fmt.Println("Database setup completed successfully.")

	// Example of creating a new note
	err = managenote.CreateNote(conn, 1, "Sample Note Title", "Sample Note Content", time.Now(), time.Now(), models.Statuses[2])
	if err != nil {
		log.Fatal("Failed to create a new note: ", err)
	}

	fmt.Println("Note created successfully.")

	// retrieving a note by ID
	note, err := managenote.GetNoteByID(conn, 1)
	if err != nil {
		log.Fatal("Failed to retrieve a note: ", err)
	}

	fmt.Printf("Note ID: %d\nTitle: %s\nNoteText: %s\nStatus: %s\nUserID: %d\n", note.NoteID, note.NoteTitle, note.NoteContent, note.Status, note.UserID)

	// updating a note
	err = managenote.UpdateNote(conn, 1, "Updated Title", "Updated Content", time.Now(), "Inactive")
	if err != nil {
		log.Fatal("Failed to update the note: ", err)
	}

	fmt.Println("Note updated successfully.")

	// deleting a note by ID
	err = managenote.DeleteNoteByID(conn, 1)
	if err != nil {
		log.Fatal("Failed to delete the note: ", err)
	}

	fmt.Println("Note deleted successfully.")
} */
