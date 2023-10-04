package managenote

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jackc/pgx/v5"
)

func TestDeleteNoteByID(t *testing.T) {
	//database connection configuration.
	connConfig, err := pgx.ParseConfig("postgres://postgres:postgres@localhost:5432/EnterpriseNotes")
	if err != nil {
		t.Fatalf("Error parsing database connection config: %v", err)
	}

	// Establish a database connection.
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	// Close the database connection when the function exits.
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

	// delete the note by its ID
	err = DeleteNoteByID(conn, 1) // Assuming 1 is the ID of the test note
	if err != nil {
		t.Fatalf("Error deleting note by ID: %v", err)
	}

	// Optionally, you can check if the note was deleted successfully by querying the database
	// and asserting that it's not present anymore.

	// Use assertions to check for errors and conditions
	assert.NoError(t, err, "DeleteNoteByID should not return an error")
}

