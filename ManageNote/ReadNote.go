package managenote

import (
	"context"
	"fmt"

	"EnterpriseNotes/models"

	"github.com/jackc/pgx/v5"
)

// GetNoteByID retrieves a note by its ID.
func GetNoteByID(conn *pgx.Conn, NoteID int) (*models.Note, error) {
	var note models.Note

	err := conn.QueryRow(context.Background(), "SELECT * FROM Notes WHERE noteID = $1", noteID).
		Scan(&note.noteID, &note.userID, &note.noteTitle, &note.noteContent, &note.creationDate, &note.completionDate, &note.status)

	if err != nil {
		// Handle the error, log it, and return it
		fmt.Println("Error retrieving note:", err)
		return nil, err
	}

	return &note, nil
}
