package models

//Create Structs for the tables
type Note struct {
	noteID         int    `db:"NoteID"`
	userID         int    `db:"UserID"`
	noteTitle      string `db:"NoteTitle"`
	noteContent    string `db:"NoteContent"`
	creationDate   string `db:"CreationDate"`
	completionDate string `db:"CompletionDate"`
	status         string `db:"Status"`
}

type User struct {
	userID           int    `db:"UserID"`
	username         string `db:"Username"`
	password         string `db:"Password"`
	email            string `db:"Email"`
	registrationDate string `db:"RegistrationDate"`
}

type Sharing struct {
	sharingID        int    `db:"SharingID"`
	noteID           int    `db:"NoteID"`
	userID           int    `db:"UserID"`
	timestamp        string `db:"Timestamp"`
	writtingSettings bool   `db:"WritingSettings"`
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
