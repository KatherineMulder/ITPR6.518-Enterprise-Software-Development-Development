package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

// Create Structs for the tables
type Note struct {
	NoteID         int       `json:"noteID"`
	UserID         int       `json:"userID"`
	NoteTitle      string    `json:"noteTitle"`
	NoteContent    string    `json:"noteContent"`
	CreationDate   time.Time `json:"creationDate"`
	DelegatedTo    string  `json:"delegatedTo"`
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
	Status string    `json:"accessLevel"` // Read or Read/Write
}

// FormattedDate formats the date for display on the web page.
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
        user_id SERIAL PRIMARY KEY, 
        username VARCHAR(100) NOT NULL UNIQUE,
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
        note_id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(user_id),
		note_title VARCHAR(50),
		note_content TEXT NOT NULL,
		creation_date TIMESTAMP NOT NULL,
		delegated_to INTEGER REFERENCES users(user_id),
		completion_date TIMESTAMP,
		status note_status
    );`

	sql = `CREATE TYPE sharing_status AS ENUM ('Read','Edit');`
	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Notes table created")


	sql = `	CREATE TYPE sharing_status AS ENUM ('Read','Edit');
	DROP TABLE IF EXISTS "sharing";
    CREATE TABLE "sharing" (
		sharing_id SERIAL PRIMARY KEY,
		note_id INTEGER,
		user_id INTEGER,
		setup_date TIMESTAMP,
		status sharing_status
	);`

	_, err = a.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sharing Table created")
	log.Printf("Inserting Data...")


	// inserting data into the "users" table//
	stmt, err := a.db.Prepare(`INSERT INTO "users"(username, password, email) VALUES($1,$2,$3)`)
	if err != nil {
		log.Fatal(err)
	}

	//readData function is used to read the data from the CSV file and store it in a slice of slices of strings.
	data, err := readData("data/Users.csv")
	if err != nil {
		log.Fatal(err)
	}


	var u User
	//range over the data slice and assign the values to the User struct.
	for _, data := range data {
		u.Username = data[0]
		u.Password = data[1]
		u.Email = data[2]

		//execute the SQL statement and pass the values from the User struct as arguments.
		_, err = stmt.Exec(u.Username, u.Password, u.Email)
		if err != nil {
			log.Fatal(err)
		}
	}



	//inserting data into the "notes" table
	stmt, err = a.db.Prepare(`INSERT INTO "notes"(user_id, note_title, note_content, creation_date, delegated_to, completion_date, status) VALUES($1,$2,$3,$4,$5,$6,$7)`)
	if err != nil {
		log.Fatal(err)
	}

	//readData function is used to read the data from the CSV file and store it in a slice of slices of strings.
	data, err = readData("data/Notes.csv")
	if err != nil {
		log.Fatal(err)
	}

	var n Note
	for _, data := range data {
		var completetime time.Time
		if data[4] != "None" {
			completetime, err = time.Parse("15:04 02-01-2006", data[3])
			if err != nil {
				log.Fatal(err)
			}
		}

		n.UserID, _ = strconv.Atoi(data[0]) //converts the string to an integer
		n.NoteTitle = data[1]
		n.NoteContent = data[2]
		//n.CreationDate = Data[3]
		//n.DelegatedTo = Data[4]  
		n.CompletionDate = completetime
		//n.Status = Data[6]

		_, err = stmt.Exec(n.UserID, n.NoteTitle, n.NoteContent, n.CompletionDate, n.Status)
		if err != nil {
			log.Fatal(err)
		}
	}



	//inserting data into the "sharing" table
	stmt, err = a.db.Prepare(`INSERT INTO sharing(note_id, user_id, setup_date, status) VALUES($1, $2, $3, $4)`)
	
	if err != nil {
		log.Fatal(err)
	}

	data, err = readData("data/Sharing.csv")
	if err != nil {
		log.Fatal(err)
	}

	///// Insertion into the "sharing" table. /////
		var s Sharing
		for _, data := range data {
			timestamp, err := time.Parse("15:04 02-01-2006", data[2])
			if err != nil {
				log.Fatal(err)
			}
			s.UserID, _ = strconv.Atoi(data[0])
			s.NoteID, _ = strconv.Atoi(data[1])
			s.Timestamp = timestamp
			s.Status = data[3] // Corrected to use the 'Status' field, not 'status'.

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
