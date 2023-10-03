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

	err := conn.QueryRow(context.Background(), "SELECT * FROM Notes WHERE noteID = $1", NoteID).
		Scan(&note.NoteID, &note.UserID, &note.NoteTitle, &note.NoteContent, &note.CreationDate, &note.CompletionDate, &note.Status)

	if err != nil {
		// Handle the error, log it, and return it
		fmt.Println("Error retrieving note:", err)
		return nil, err
	}

	return &note, nil
}
