package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
	//"golang.org/x/crypto/bcrypt" // Importing the bcrypt package for password hashing
)

// Create Structs for the tables
type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Note struct {
	NoteID         int       `json:"noteID"`
	UserID         int       `json:"userID"`
	NoteTitle      string    `json:"noteTitle"`
	NoteContent    string    `json:"noteContent"`
	CreationDate   time.Time `json:"creationDate"`
	DelegatedTo    string    `json:"delegatedTo"`
	CompletionDate time.Time `json:"completionDate"`
	Status         string    `json:"status"`
	Privileges     string
	SharedUsers    []Sharing
}

type Sharing struct {
	SharingID int       `json:"sharingID"`
	NoteID    int       `json:"noteID"`
	UserID    int       `json:"userID"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

func (n Note) FormattedDate() string {
	return n.CompletionDate.Format(time.ANSIC)
}

// Read data from csv file
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

// //////creating tables in a PostgreSQL database and inserting data into them.///////
func (a *App) importData() error {
	log.Printf("Creating tables...")

	sql := `DROP TABLE IF EXISTS "users";
	CREATE TABLE "users" (
		userID SERIAL PRIMARY KEY, 
		username VARCHAR(100) NOT NULL,
		password VARCHAR(100) NOT NULL
	);
	CREATE UNIQUE INDEX users_by_id ON "users" (userID);`

	_, err := a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Users table created")

	sql = `CREATE TYPE note_status AS ENUM ('None','In Progress','Completed','Cancelled');
	DROP TABLE IF EXISTS "notes";
    CREATE TABLE "notes" (
        noteID SERIAL PRIMARY KEY,
        userID INTEGER NOT NULL,
		note_title VARCHAR(50),
		note_content TEXT NOT NULL,
		creationDate TIMESTAMP NOT NULL,
		delegatedTo VARCHAR(100),
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

	// inserting data into the "users" table
	stmt, err := a.db.Prepare("INSERT INTO users VALUES($1,$2,$3)")
	if err != nil {
		log.Fatal(err)
	}

	//readData function is used to read the data from the CSV file and store it in a slice of slices of strings.
	data, err := readData("data/Users.csv")
	if err != nil {
		log.Fatal(err)
	}
	/*// Prepare the user_shares insert query this is an example code for how to handle the password hash ....
	userSharesStmt, err := a.db.Prepare("INSERT INTO user_shares (note_id, username, privileges) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}*/
	var u User
	//range over the data slice and assign the values to the User struct.
	for _, data := range data {
		u.UserID, _ = strconv.Atoi(data[0])
		u.Username = data[1]
		u.Password = data[2]
		_, err := stmt.Exec(u.UserID, u.Username, u.Password)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Inserted Data to usersTable")

	//inserting data into the "notes" table
	stmt, err = a.db.Prepare(`INSERT INTO "notes"(userID, note_title, note_content, creationDate, delegatedTo, completion_date, status) VALUES($1,$2,$3,$4,$5,$6,$7)`)
	if err != nil {
		log.Fatal(err)
	}

	//readData function is used to read the data from the CSV file and store it in a slice of slices of strings.
	data, err = readData("data/Notes.csv")
	if err != nil {
		log.Fatal(err)
	}

	/////insertion into the "notes" table.///////
	var n Note
	for _, data := range data {
		/*var completetime *time.Time
		creationDate := time.Now()

		if data[3] != "None" {
			parsedTime, err := time.Parse("02/01", data[3])
			if err != nil {
				parsedTime, err = time.Parse("2006-02-02", data[3])
				if err != nil {
					log.Printf("Error parsing date for row: %v, error: %v", data, err)
					continue
				}
			}
			creationDate = parsedTime
		}*/

		n.UserID, _ = strconv.Atoi(data[0])
		n.NoteTitle = data[1]
		n.NoteContent = data[2]
		n.CreationDate = time.Now()
		n.DelegatedTo = data[4]
		n.CompletionDate = time.Now()
		n.Status = data[6]

		_, err = stmt.Exec(n.UserID, n.NoteTitle, n.NoteContent, n.CreationDate, n.DelegatedTo, n.CompletionDate, n.Status)
		if err != nil {
			log.Fatal(err)
		}

	}
	log.Printf("Inserted Data to notesTable")

	//inserting data into the "sharing" table
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
		timestamp, err := time.Parse("15:04 02-01-2006", data[2])
		if err != nil {
			log.Fatal(err)
		}
		s.UserID, _ = strconv.Atoi(data[0])
		s.NoteID, _ = strconv.Atoi(data[1])
		s.Timestamp = timestamp
		s.Status = data[3]

		_, err = stmt.Exec(s.UserID, s.NoteID, s.Timestamp, s.Status)
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
