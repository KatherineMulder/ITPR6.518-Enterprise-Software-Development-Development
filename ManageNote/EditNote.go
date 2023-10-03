package managenote

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// UpdateNote updates an existing note.
func UpdateNote(conn *pgx.Conn, NoteID int, NoteTitle string, NoteContent string, CompletionDate time.Time, Status string) error {
	_, err := conn.Exec(context.Background(), "UPDATE Notes SET noteTitle = $1, noteContent = $2, completionDate = $3, status = $4 WHERE noteID = $6", NoteTitle, NoteContent, CompletionDate, Status[3], NoteID)
	if err != nil {
		return err
	}

	return nil
}
