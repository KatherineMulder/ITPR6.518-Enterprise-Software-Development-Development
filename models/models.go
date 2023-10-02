package models

import "time"

//Create Structs for the tables
type Note struct {
	ID               int       `db:"noteID"`
	UserID           int       `db:"userID"`
	NoteTitle         string    `db:"noteTitle"`
	NoteContent         string    `db:"NoteContent"`
	CreationDateTime time.Time `db:"creationDateTime"`
	CompletionDateTime time.Time `db:"completionDateTime"`
	Status           string    `db:"status"`
	DelegatedToUserID int       `db:"delegatedToUserID"`
}

type User struct {
	ID              int       `db:"userID"`
	Username        string    `db:"name"`
	Password        string    `db:"password_hash"`
	Email           string    `db:"email"`
	RegistrationDate time.Time `db:"registration_date"`
}


type Sharing struct {
	SharingID int       `db:"sharingID"`
	NoteID    int       `db:"noteID"`
	UserID    int       `db:"userID"`
	Status    string    `db:"status"`
	Timestamp time.Time `db:"timestamp"`
}

// CreateStructsForTestData creates and returns sample data instances for testing.
func CreateStructsForTestData() (Note, User, Sharing) {
	return Note{
			ID:               1,
			UserID:           1,
			NoteTitle:         "Sample Note",
			NoteContent:         "This is a sample note content.",
			CreationDateTime: time.Now(),
			CompletionDateTime: time.Now().Add(24 * time.Hour), // completion time
			Status:           "In Progress",
			DelegatedToUserID: 2, // User ID to whom it's delegated
		},
		User{
			ID:              1,
			Username:        "exampleUser",
			Password:        "hashedPassword", 
			Email:           "user@example.com",
			RegistrationDate: time.Now(),
		},
		Sharing{
			SharingID: 1,
			NoteID:    1,
			UserID:    3, // User ID with whom the note is shared
			Status:    "Read",
			Timestamp: time.Now(),
		}
}

