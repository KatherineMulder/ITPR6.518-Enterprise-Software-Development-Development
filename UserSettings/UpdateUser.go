package usersettings

import (
	"context"
	

	"github.com/jackc/pgx/v5"
)

func UpdateUser(conn *pgx.Conn, UserID int, Username string, Password string, Email string) error {
	_, err := conn.Exec(context.Background(), "UPDATE Users SET username = $1, password = $2, email = $3 WHERE userID = $4", Username, Password, Email, UserID)

	if err != nil {
		return err
	}

	return nil
}
