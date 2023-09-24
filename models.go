package main

import "time"

//Create Structs
type Note struct {
    ID        int
    Title     string
    Content   string
    Timestamp time.Time
    Status    string
    UserID    int
}

type User struct {
	ID       int
	Username string
	Password string
	Email    string
}

type Sharing struct {
	sharingID int
	NoteID    int
	UserID    int
	status    string
	Timestamp time.Time
}

// Create Struct for Test Data


