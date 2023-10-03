package managenote

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// CreateNote inserts a new note into the Notes table.
func CreateNote(conn *pgx.Conn, UserID int, NoteTitle string, NoteContent string, CreationDate time.Time, CompletionDate time.Time, Status string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO Notes (userID, noteTitle, noteContent, creationDate, completiondate, status) VALUES ($1, $2, $3, $4, $5, $6)",
		UserID, NoteTitle, NoteContent, CompletionDate, CompletionDate, Status[2])

	if err != nil {
		return err
	}

	return nil
}
