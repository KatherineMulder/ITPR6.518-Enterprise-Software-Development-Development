package usersettings

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func CreateUser(conn *pgx.Conn, Username string, Password string, Email string, RegistrationDate time.Time) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO Users(username, password, email, registrationDate) VALUES ($1, $2, $3, $4)", Username, Password, Email, RegistrationDate)

	if err != nil {
		return err
	}

	return nil
}
