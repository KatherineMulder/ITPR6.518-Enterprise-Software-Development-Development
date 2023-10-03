package databasesetup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

var databaseURL = "postgres://postgres:postgres@localhost:5432/EnterpriseNotes"

func DatabaseSetup() (*pgx.Conn, error) {

	// Use the databaseURL variable for the connection string
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		// If the database doesn't exist, create it
		if err := createDatabase(ctx); err != nil {
			log.Fatal("Failed to create the database: ", err)
			return nil, err
		}

		// Reconnect to the newly created database
		conn.Close(ctx)
		conn, err = pgx.Connect(context.Background(), databaseURL)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	fmt.Println("Connected successfully")

	if err := createTables(conn); err != nil {
		log.Fatal("Failed to create tables: ", err)
		return nil, err
	}

	return conn, nil
}

func createDatabase(ctx context.Context) error {
	// Connect to PostgreSQL without specifying a database
	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// Check if the database already exists
	var dbExists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", "EnterpriseNotes").Scan(&dbExists)
	if err != nil {
		return err
	}

	if !dbExists {
		// Create the database if it doesn't exist
		_, err := conn.Exec(ctx, "CREATE DATABASE EnterpriseNotes")
		if err != nil {
			return err
		}
		fmt.Println("Database created successfully")
	}

	return nil
}

// Create Tables Function
func createTables(conn *pgx.Conn) error {

	usersTable := `DROP TABLE IF EXISTS "user";
    CREATE TABLE "user" (
        userID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, 
        username VARCHAR(100),
		password VARCHAR(100),
		email VARCHAR(100),
		registrationDate DATE DEFAULT CURRENT_TIMESTAMP
    );
    `
	notesTable := `DROP TABLE IF EXISTS notes;
    CREATE TABLE Notes (
        noteID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        userID INT,
		noteTitle VARCHAR(50),
		noteContent TEXT,
		creationDate timestamp DEFAULT CURRENT_TIMESTAMP,
		completionDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(50)
    );
    `

	sharingTable := `DROP TABLE IF EXISTS sharing;
    CREATE TABLE Sharing (
		sharingID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		noteID INT,
		userID INT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		writingSetting BOOL DEFAULT false
	);
    `
	_, err := conn.Exec(context.Background(), usersTable)
	if err != nil {
		log.Fatalf("An error occurred when creating the 'users' table.\nGot %s", err)
		return err
	}

	_, err = conn.Exec(context.Background(), notesTable)
	if err != nil {
		log.Fatalf("An error occurred when creating the 'notes' table.\nGot %s\n", err)
		return err
	}

	_, err = conn.Exec(context.Background(), sharingTable)
	if err != nil {
		log.Fatalf("An error occurred when creating the 'sharings' table.\nGot %s\n", err)
		return err
	}

	// No error occurred, so return nil to indicate success
	return nil
}
