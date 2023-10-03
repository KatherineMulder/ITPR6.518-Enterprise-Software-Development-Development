package managenote

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// DeleteNoteByID deletes a note by its ID.
func DeleteNoteByID(conn *pgx.Conn, NoteID int) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM Notes WHERE noteID = $1", NoteID)
	if err != nil {
		return err
	}

	return nil
}
