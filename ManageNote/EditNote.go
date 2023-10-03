package managenote

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// UpdateNote updates an existing note.
func UpdateNote(conn *pgx.Conn, noteID int, noteTitle string, noteContent string, completionDate time.Time, status string) error {
	_, err := conn.Exec(context.Background(), "UPDATE Notes SET NoteTitle = $1, NoteContent = $2, CompletionDate = $3, Status = $4 WHERE noteID = $6", noteTitle, noteContent, completionDate, status[3], noteID)
	if err != nil {
		return err
	}

	return nil
}
