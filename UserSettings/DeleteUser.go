package usersettings

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func DeleteUserByID(conn *pgx.Conn, UserID int) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM Users WHERE userID = $1", UserID)

	if err != nil {
		return err
	}

	return nil
}
