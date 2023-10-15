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
	CompletionDate time.Time `json:"completionDate"`
	Status         string    `json:"status"`
}

type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Sharing struct {
	SharingID int       `json:"sharingID"`
	NoteID    int       `json:"noteID"`
	UserID    int       `json:"userID"`
	Timestamp time.Time `json:"timestamp"`
	status    string    `json:"status"`
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

	sql := `DROP TABLE IF EXISTS "users";
    CREATE TABLE "users" (
        userID INTEGER PRIMARY KEY NOT NULL, 
        username VARCHAR(100) NOT NULL,
		password VARCHAR(100) NOT NULL,
		email VARCHAR(100)
    );`
	_, err := a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Users table created")

	sql = `DROP TABLE IF EXISTS "notes";
    CREATE TABLE "notes" (
        noteID INTEGER PRIMARY KEY NOT NULL,
        userID INTEGER NOT NULL,
		note_title VARCHAR(50),
		note_content TEXT NOT NULL,
		completion_date TIMESTAMP,
		status VARCHAR(50)
    );`
	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Notes table created")

	sql = `	CREATE TYPE sharing_status AS ENUM ('Read','Edit');
	DROP TABLE IF EXISTS "sharing";
    CREATE TABLE "sharing" (
		sharingID INTEGER PRIMARY KEY NOT NULL,
		noteID INTEGER,
		userID INTEGER,
		setup_date TIMESTAMP,
		status sharing_status
	);`
	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sharing Table created")

	log.Printf("Inserting Data...")

	stmt, err := a.db.Prepare(`INSERT INTO "users" VALUES($1,$2,$3,$4)`)
	if err != nil {
		log.Fatal(err)
	}

	data, err := readData("data/Users.csv")
	if err != nil {
		log.Fatal(err)
	}

	var u User
	for _, data := range data {
		u.UserID, _ = strconv.Atoi(data[0])
		u.Username = data[1]
		u.Password = data[2]
		u.Email = data[3]

		_, err = stmt.Exec(u.UserID, u.Username, u.Password, u.Email)
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt, err = a.db.Prepare(`INSERT INTO "notes" VALUES($1,$2,$3,$4,$5,$6)`)
	if err != nil {
		log.Fatal(err)
	}

	data, err = readData("data/Notes.csv")
	if err != nil {
		log.Fatal(err)
	}

	var n Note
	for _, data := range data {
		completetime, err := time.Parse("2006-01-02 15:04", data[4])
		if err != nil {
			log.Fatal(err)
		}
		n.NoteID, _ = strconv.Atoi(data[0])
		n.UserID, _ = strconv.Atoi(data[1])
		n.NoteTitle = data[2]
		n.NoteContent = data[3]
		n.CompletionDate = completetime
		n.Status = data[5]

		_, err = stmt.Exec(n.NoteID, n.UserID, n.NoteTitle, n.NoteContent, n.CompletionDate, n.Status)
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt, err = a.db.Prepare(`INSERT INTO "sharing" VALUES($1,$2,$3,$4,$5)`)
	if err != nil {
		log.Fatal(err)
	}

	data, err = readData("data/Sharing.csv")
	if err != nil {
		log.Fatal(err)
	}

	var s Sharing
	for _, data := range data {
		timestamp, err := time.Parse("2006-01-02 15:04", data[3])
		if err != nil {
			log.Fatal(err)
		}
		s.SharingID, _ = strconv.Atoi(data[0])
		s.UserID, _ = strconv.Atoi(data[1])
		s.NoteID, _ = strconv.Atoi(data[2])
		s.Timestamp = timestamp
		s.status = data[4]

		_, err = stmt.Exec(s.SharingID, s.UserID, s.NoteID, s.Timestamp, s.status)
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
