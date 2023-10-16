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
        userID SERIAL PRIMARY KEY, 
        username VARCHAR(100) NOT NULL,
		password VARCHAR(100) NOT NULL,
		email VARCHAR(100)
    );`
	_, err := a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Users table created")

	sql = `CREATE TYPE note_status AS ENUM ('None','In Progress','Completed','Cancelled','Delegated');
	DROP TABLE IF EXISTS "notes";
    CREATE TABLE "notes" (
        noteID SERIAL PRIMARY KEY,
        userID INTEGER NOT NULL,
		note_title VARCHAR(50),
		note_content TEXT NOT NULL,
		completion_date TIMESTAMP,
		status note_status
    );`
	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Notes table created")

	sql = `	CREATE TYPE sharing_status AS ENUM ('Read','Edit');
	DROP TABLE IF EXISTS "sharing";
    CREATE TABLE "sharing" (
		sharingID SERIAL PRIMARY KEY,
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

	stmt, err := a.db.Prepare(`INSERT INTO "users"(username, password, email) VALUES($1,$2,$3)`)
	if err != nil {
		log.Fatal(err)
	}

	data, err := readData("data/Users.csv")
	if err != nil {
		log.Fatal(err)
	}

	var u User
	for _, data := range data {
		u.Username = data[0]
		u.Password = data[1]
		u.Email = data[2]

		_, err = stmt.Exec(u.Username, u.Password, u.Email)
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt, err = a.db.Prepare(`INSERT INTO "notes"(userID, note_title, note_content, completion_date, status) VALUES($1,$2,$3,$4,$5)`)
	if err != nil {
		log.Fatal(err)
	}

	data, err = readData("data/Notes.csv")
	if err != nil {
		log.Fatal(err)
	}

	var n Note
	for _, data := range data {
		completetime := time.Now()
		if data[4] != "None" {
			completetime, err = time.Parse("2006-01-02 15:04", data[3])
			if err != nil {
				log.Fatal(err)
			}

		}
		n.UserID, _ = strconv.Atoi(data[0])
		n.NoteTitle = data[1]
		n.NoteContent = data[2]
		n.CompletionDate = completetime
		n.Status = data[4]

		_, err = stmt.Exec(n.UserID, n.NoteTitle, n.NoteContent, n.CompletionDate, n.Status)
		if err != nil {
			log.Fatal(err)
		}
	}

	stmt, err = a.db.Prepare(`INSERT INTO "sharing"(noteID, userID, setup_date, status) VALUES($1,$2,$3,$4)`)
	if err != nil {
		log.Fatal(err)
	}

	data, err = readData("data/Sharing.csv")
	if err != nil {
		log.Fatal(err)
	}

	var s Sharing
	for _, data := range data {
		timestamp, err := time.Parse("2006-01-02 15:04", data[2])
		if err != nil {
			log.Fatal(err)
		}
		s.UserID, _ = strconv.Atoi(data[0])
		s.NoteID, _ = strconv.Atoi(data[1])
		s.Timestamp = timestamp
		s.status = data[3]

		_, err = stmt.Exec(s.UserID, s.NoteID, s.Timestamp, s.status)
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
