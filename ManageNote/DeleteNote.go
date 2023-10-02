package managenote

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// DeleteNoteByID deletes a note by its ID.
func DeleteNoteByID(conn *pgx.Conn, noteID int) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM Notes WHERE noteID = $1", noteID)
	if err != nil {
		return err
	}

	return nil
}
