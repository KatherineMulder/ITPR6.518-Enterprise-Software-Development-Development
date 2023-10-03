package models

//Create Structs for the tables
type Note struct {
	NoteID         int    `db:"noteID"`
	UserID         int    `db:"userID"`
	NoteTitle      string `db:"noteTitle"`
	NoteContent    string `db:"noteContent"`
	CreationDate   string `db:"creationDate"`
	CompletionDate string `db:"completionDate"`
	Status         string `db:"status"`
}

type User struct {
	UserID           int    `db:"userID"`
	Username         string `db:"username"`
	Password         string `db:"password"`
	Email            string `db:"email"`
	RegistrationDate string `db:"registrationDate"`
}

type Sharing struct {
	SharingID        int    `db:"sharingID"`
	NoteID           int    `db:"noteID"`
	UserID           int    `db:"userID"`
	Timestamp        string `db:"timestamp"`
	WrittingSettings bool   `db:"writingSettings"`
}

// Create Slices
var Users []User
var Notes []Note
var Sharings []Sharing
var Statuses = [5]string{
	"none",
	"in progress",
	"completed",
	"cancelled",
	"delegated",
}
