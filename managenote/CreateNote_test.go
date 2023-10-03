package managenote

import (
    "testing"
    "context"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/jackc/pgx/v5"
)

func TestCreateNote(t *testing.T) {
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
    // Close the database connection when the function exits.
	defer conn.Close(context.Background())


    // Define test data
    userID := 1
    noteTitle := "Test use Note"
    noteContent := "This is a test note."
    creationDate := time.Now()
    completionDate := time.Now()
    status := "Pending"

    // Call the function we are testing
    err = CreateNote(conn, userID, noteTitle, noteContent, creationDate, completionDate, status)

   // Use assertions to check for errors and conditions
   assert.NoError(t, err, "CreateNote should not return an error")
}
