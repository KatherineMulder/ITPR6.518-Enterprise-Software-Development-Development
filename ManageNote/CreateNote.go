package managenote

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// CreateNote inserts a new note into the Notes table.
func CreateNote(conn *pgx.Conn, userID int, noteName string, noteText string, creationDateTime time.Time, completionDateTime time.Time, status string, delegatedToUserID int) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO Notes (userid, noteName, notetext, creationdatetime, completiondatetime, status, delegatedtouserid) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		userID, noteName, noteText, creationDateTime, completionDateTime, status, delegatedToUserID)

	if err != nil {
		return err
	}

	return nil
}
