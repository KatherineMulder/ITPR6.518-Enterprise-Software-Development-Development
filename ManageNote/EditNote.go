package managenote

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// UpdateNote updates an existing note.
func UpdateNote(conn *pgx.Conn, noteID int, noteName string, noteText string, completionDateTime time.Time, status string, delegatedToUserID int) error {
	_, err := conn.Exec(context.Background(), "UPDATE Notes SET noteName = $1, noteText = $2, completionDateTime = $3, status = $4, delegatedToUserID = $5 WHERE noteID = $6", noteName, noteText, completionDateTime, status, delegatedToUserID, noteID)
	if err != nil {
		return err
	}

	return nil
}
