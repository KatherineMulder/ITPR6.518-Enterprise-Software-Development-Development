package managenote

/*import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jackc/pgx/v5"
)

func TestGetNoteByID(t *testing.T) {
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
	//lose the database when the function exits.
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
	if err != nil {
		t.Fatalf("Error inserting test note: %v", err)
	}

	//retrieve the note by its ID
	retrievedNote, err := GetNoteByID(conn, 1) // Assuming 1 is the ID of the test note
	if err != nil {
		t.Fatalf("Error retrieving note by ID: %v", err)
	}

	//check if the retrieved note matches the expected values.
	assert.Equal(t, userID, retrievedNote.UserID, "User ID should match")
	assert.Equal(t, noteTitle, retrievedNote.NoteTitle, "Note title should match")
	assert.Equal(t, noteContent, retrievedNote.NoteContent, "Note content should match")
	assert.Equal(t, status, retrievedNote.Status, "Status should match")

	//check for errors and conditions
	assert.NoError(t, err, "GetNoteByID should not return an error")
}
*/
