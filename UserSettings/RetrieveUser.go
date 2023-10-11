package usersettings

import (
	"EnterpriseNotes/models"
	"context"

	"github.com/jackc/pgx/v5"
)

func RetrieveUserByID(conn *pgx.Conn, UserID int) (*models.User, error) { // Use models.User as the return type
	var user models.User // Use models.User as the variable type

	row := conn.QueryRow(context.Background(), "SELECT userID, username, password, email, registrationDate FROM Users WHERE userID = $1", UserID)

	err := row.Scan(&user.UserID, &user.Username, &user.Password, &user.Email, &user.RegistrationDate)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Handle the case where no user with the given UserID is found
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
