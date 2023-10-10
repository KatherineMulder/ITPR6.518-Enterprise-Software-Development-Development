package models

import
"time"

//Create Structs for the tables
type Note struct {
	NoteID         int    `json:"noteID"`
	UserID         int    `json:"userID"`
	NoteTitle      string `json:"noteTitle"`
	NoteContent    string `json:"noteContent"`
	CreationDate   time.Time `json:"creationDate"`
	CompletionDate time.Time `json:"completionDate"`
	Status         string `json:"status"`
}

type User struct {
	UserID           int    `json:"userID"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Email            string `json:"email"`
	RegistrationDate string `json:"registrationDate"`
}

type Sharing struct {
	SharingID        int    `json:"sharingID"`
	NoteID           int    `json:"noteID"`
	UserID           int    `json:"userID"`
	Timestamp        time.Time `json:"timestamp"`
	WrittingSettings bool   `json:"writingSettings"`
}

// Create Struct for Test Data
type Data struct {
	Users        []User        `json:"users"`
	Notes        []Note        `json:"notes"`
	Sharings     []Sharing     `json:"sharings"`
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

