package models

//Create Structs for the tables
type Note struct {
	ID               	int       `db:"noteID"`
	UserID           	int       `db:"userID"`
	DelegatedToUserID 	int       `db:"delegatedToUserID"`
	NoteTitle         	string    `db:"noteTitle"`
	NoteContent         string    `db:"NoteContent"`
	CreationDate 		string 	  `db:"creationDateTime"`
	CompletionDate  	string 	  `db:"completionDateTime"`
	Status           	string    `db:"status"`
	sharedUsers			string    `db:"sharedUsers"`
	
}

type User struct {
	ID              	int       `db:"userID"`
	Username        	string    `db:"name"`
	Password        	string    `db:"password_hash"`
	Email           	string    `db:"email"`
	RegistrationDate 	string	  `db:"registration_date"`
	readingSettings		bool    `db:"readingSettings"`
	writtingSettings	bool    `db:"writtingSettings"`
}


type Sharing struct {
	SharingID 			int       `db:"sharingID"`
	NoteID    			int       `db:"noteID"`
	UserID   			int       `db:"userID"`
	Status    			string    `db:"status"`
	Timestamp 			string 	  `db:"timestamp"`
}

// Create Struct for Test Data
type Data struct {
	Users        []User        `db:"users"`
	Notes        []Note        `db:"notes"`
	Sharings     []Sharing     `db:"associations"`
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
