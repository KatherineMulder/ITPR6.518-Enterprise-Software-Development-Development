package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

// Create Structs for the tables
type Note struct {
	NoteID         int       `json:"noteID"`
	UserID         int       `json:"userID"`
	NoteTitle      string    `json:"noteTitle"`
	NoteContent    string    `json:"noteContent"`
	CreationDate   time.Time `json:"creationDate"`
	CompletionDate time.Time `json:"completionDate"`
	Status         string    `json:"status"`
}

type User struct {
	UserID           int       `json:"userID"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	Email            string    `json:"email"`
	Role             string    `json:"role"`
	RegistrationDate time.Time `json:"registrationDate"`
}

type Sharing struct {
	SharingID        int       `json:"sharingID"`
	NoteID           int       `json:"noteID"`
	UserID           int       `json:"userID"`
	Timestamp        time.Time `json:"timestamp"`
	WrittingSettings bool      `json:"writingSettings"`
}

func readData(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func (a *App) importData() error {
	log.Printf("Creating tables...")

	sql := `DROP TABLE IF EXISTS user;
    CREATE TABLE User (
        userID INTEGER PRIMARY KEY NOT NULL, 
        username VARCHAR(100) NOT NULL,
		password VARCHAR(100) NOT NULL,
		email VARCHAR(100),
		role INTEGER DEFAULT 2 NOT NULL,
		registrationDate DATE DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("User table created")

	sql = `DROP TABLE IF EXISTS notes;
    CREATE TABLE Notes (
        noteID INTEGER PRIMARY KEY NOT NULL,
        userID INTEGER NOT NULL,
		noteTitle VARCHAR(50) NOT NULL,
		noteContent TEXT NOT NULL,
		creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completionDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(50) DEFAULT none
    );`
	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Notes tble created")

	sql = `DROP TABLE IF EXISTS sharing;
    CREATE TABLE Sharing (
		sharingID INTEGER PRIMARY KEY NOT NULL,
		noteID INTEGER,
		userID INTEGER,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		writingSetting BOOL DEFAULT false
	);`
	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sharing Table created")

	log.Printf("Inserting Data...")

	stmt, err := a.db.Prepare("INSERT INTO Users VALUES($1,$2,$3,$4,$5,$6)")
	if err != nil {
		log.Fatal(err)
	}

	data, err := readData("data/Users.csv")
	if err != nil {
		log.Fatal(err)
	}

	var u User
	for _, data := range data {
		registrationtime, err := time.Parse("2006-01-02 15:04", data[6])
		if err != nil {
			log.Fatal(err)
		}
		u.UserID, _ = strconv.Atoi(data[1])
		u.Username = data[2]
		u.Password = data[3]
		u.Email = data[4]
		u.Role = data[5]
		u.RegistrationDate = registrationtime

		_, err = stmt.Exec(u.UserID, u.Username, u.Password, u.Email, u.Role, u.RegistrationDate)
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt, err = a.db.Prepare("INSERT INTO Notes VALUES($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		log.Fatal(err)
	}

	data, err = readData("data/Notes.csv")
	if err != nil {
		log.Fatal(err)
	}

	var n Note
	for _, data := range data {
		completetime, err := time.Parse("2006-01-02 15:04", data[5])
		if err != nil {
			log.Fatal(err)
		}
		creationtime, err := time.Parse("2006-01-02 15:04", data[6])
		if err != nil {
			log.Fatal(err)
		}
		n.NoteID, _ = strconv.Atoi(data[1])
		n.UserID, _ = strconv.Atoi(data[2])
		n.NoteTitle = data[3]
		n.NoteContent = data[4]
		n.CompletionDate = completetime
		n.CreationDate = creationtime
		n.Status = data[7]

		_, err = stmt.Exec(n.NoteID, n.UserID, n.NoteTitle, n.NoteContent, n.CompletionDate, n.CreationDate, n.Status)
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt, err = a.db.Prepare("INSERT INTO Sharing VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		log.Fatal(err)
	}

	data, err = readData("data/Sharing.csv")
	if err != nil {
		log.Fatal(err)
	}

	var s Sharing
	for _, data := range data {
		timestamp, err := time.Parse("2006-01-02 15:04", data[4])
		if err != nil {
			log.Fatal(err)
		}
		writingsettings, err := strconv.ParseBool(data[5])
		if err != nil {
			log.Fatal(err)
		}
		s.SharingID, _ = strconv.Atoi(data[1])
		s.UserID, _ = strconv.Atoi(data[2])
		s.NoteID, _ = strconv.Atoi(data[3])
		s.Timestamp = timestamp
		s.WrittingSettings = writingsettings

		_, err = stmt.Exec(s.SharingID, s.UserID, s.NoteID, s.Timestamp, s.WrittingSettings)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Create("./imported")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return err
}
