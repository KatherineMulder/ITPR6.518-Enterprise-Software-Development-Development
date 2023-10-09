package managenote

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jackc/pgx/v5"
)

func TestUpdateNote(t *testing.T) {
	// Define the database connection configuration.
	connConfig, err := pgx.ParseConfig("postgres://postgres:postgres@localhost:5432/EnterpriseNotes")
	if err != nil {
		t.Fatalf("Error parsing database connection config: %v", err)
	}

	// Establish a database connection.
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	// Close the database when the function exits.
	defer conn.Close(context.Background())

	// Define test data
	// Insert a test note into the database first
	userID := 1
	noteTitle := "Test use Note"
	noteContent := "This is a test note."
	creationDate := time.Now()
	completionDate := time.Now()
	status := "Pending"

	err = CreateNote(conn, userID, noteTitle, noteContent, creationDate, completionDate, status)
	if err != nil {
		t.Fatalf("Error inserting test note: %v", err)
	}

	// Update the note by its ID
	newNoteTitle := "Updated Note Title"
	newNoteContent := "Updated note content."
	newCompletionDate := time.Now().Add(24 * time.Hour)
	newStatus := "Completed"

	err = UpdateNote(conn, 1, newNoteTitle, newNoteContent, newCompletionDate, newStatus) // Assuming 1 is the ID of the test note
	if err != nil {
		t.Fatalf("Error updating note: %v", err)
	}

	// Query the database to retrieve the updated note
	//retrievedNote, err := GetNoteByID(conn, 1) // Assuming 1 is the ID of the test note
	//if err != nil {
		//t.Fatalf("Error retrieving updated note: %v", err)
	//}

	//retrieved values match the updated ones
	//assert.Equal(t, newNoteTitle, retrievedNote.NoteTitle, "Note title should match")
	//assert.Equal(t, newNoteContent, retrievedNote.NoteContent, "Note content should match")
	//assert.Equal(t, newCompletionDate, retrievedNote.CompletionDate, "Completion date should match")
	//assert.Equal(t, newStatus, retrievedNote.Status, "Status should match")

	// Use assertions to check for errors and conditions
	assert.NoError(t, err, "UpdateNote should not return an error")
}
