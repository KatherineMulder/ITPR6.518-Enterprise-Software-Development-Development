package managenote

import (
    "testing"
    "context"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/jackc/pgx/v5"
)

func TestCreateNote(t *testing.T) {
    //database connection configuration.
    connConfig, err := pgx.ParseConfig("postgres://postgres:postgres@localhost:5432/EnterpriseNotes")
    if err != nil {
        t.Fatalf("Error parsing database connection config: %v", err)
    }

    //database connection.
    conn, err := pgx.ConnectConfig(context.Background(), connConfig)
    if err != nil {
        t.Fatalf("Failed to connect to the database: %v", err)
    }
    //close the database when the function exits.
	defer conn.Close(context.Background())


    //test data
    userID := 1
    noteTitle := "Test use Note"
    noteContent := "This is a test note."
    creationDate := time.Now()
    completionDate := time.Now()
    status := "Pending"

    //insert a test note into the database
    err = CreateNote(conn, userID, noteTitle, noteContent, creationDate, completionDate, status)

   //check for errors and conditions
   assert.NoError(t, err, "CreateNote should not return an error")
}
