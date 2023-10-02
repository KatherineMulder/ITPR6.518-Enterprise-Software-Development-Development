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

func createTables(conn *pgx.Conn) error {
	usersTable := `
    CREATE TABLE IF NOT EXISTS Users (
        userID serial PRIMARY KEY,
        name VARCHAR(255),
        email VARCHAR(255) UNIQUE,
        password_hash VARCHAR(255),
        registration_date TIMESTAMP
    );
    `
	notesTable := `
    CREATE TABLE IF NOT EXISTS Notes (
        noteID serial PRIMARY KEY,
        userID INT,
        noteTitle VARCHAR(255),
        NoteContent TEXT,
        creationDateTime TIMESTAMP,
        completionDateTime TIMESTAMP,
        status VARCHAR(255),
        delegatedToUserID INT
    );
    `

	sharingTable := `
    CREATE TABLE IF NOT EXISTS Sharing (
        sharingID serial PRIMARY KEY,
        noteID INT,
        userID INT,
        status VARCHAR(255),
        timestamp TIMESTAMP
    );
    `

	_, err := conn.Exec(context.Background(), usersTable)
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), notesTable)
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), sharingTable)
	if err != nil {
		return err
	}

	return nil
}

// DeleteNoteByID deletes a note by its ID.
func DeleteNoteByID(conn *pgx.Conn, noteID int) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM Notes WHERE noteID = $1", noteID)
	if err != nil {
		return err
	}

	return nil
}
