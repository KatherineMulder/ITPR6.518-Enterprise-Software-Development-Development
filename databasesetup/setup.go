package databasesetup

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/jackc/pgx/v5"
)

var databaseURL = "postgres://postgres:postgres@localhost:5432/EnterpriseNotes"


type Note struct {
    ID        int
    Title     string
    Content   string
    Timestamp time.Time
    Status    string
    UserID    int
}

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
    notesTable := `
    CREATE TABLE IF NOT EXISTS Notes (
        ID serial PRIMARY KEY,
        title VARCHAR(255),
        content TEXT,
        timestamp TIMESTAMP,
        status VARCHAR(255),
        userID INT
    );
    `

    usersTable := `
    CREATE TABLE IF NOT EXISTS Users (
        userID serial PRIMARY KEY,
        name VARCHAR(255),
        email VARCHAR(255) UNIQUE,
        password_hash VARCHAR(255)
    );
    `

    sharingTable := `
    CREATE TABLE IF NOT EXISTS Sharing (
        userID INT,
        noteID INT,
        status VARCHAR(255),
        timestamp TIMESTAMP
    );
    `

    _, err := conn.Exec(context.Background(), notesTable)
    if err != nil {
        return err
    }

    _, err = conn.Exec(context.Background(), usersTable)
    if err != nil {
        return err
    }

    _, err = conn.Exec(context.Background(), sharingTable)
    if err != nil {
        return err
    }

    return nil
}

// CreateNote inserts a new note into the Notes table.
func CreateNote(conn *pgx.Conn, title string, content string, status string, userID int) error {
    _, err := conn.Exec(context.Background(), "INSERT INTO Notes (title, content, timestamp, status, userID) VALUES ($1, $2, $3, $4, $5)",
        title, content, time.Now(), status, userID)
    if err != nil {
        return err
    }

    return nil
}

// GetNoteByID retrieves a note by its ID.
func GetNoteByID(conn *pgx.Conn, noteID int) (*Note, error) {
    var note Note

    err := conn.QueryRow(context.Background(), "SELECT * FROM Notes WHERE ID = $1", noteID).
        Scan(&note.ID, &note.Title, &note.Content, &note.Timestamp, &note.Status, &note.UserID)
    if err != nil {
        return nil, err
    }

    return &note, nil
}

// UpdateNote updates an existing note.
func UpdateNote(conn *pgx.Conn, noteID int, title string, content string, status string) error {
    _, err := conn.Exec(context.Background(), "UPDATE Notes SET title = $1, content = $2, status = $3 WHERE ID = $4", title, content, status, noteID)
    if err != nil {
        return err
    }

    return nil
}

// DeleteNoteByID deletes a note by its ID.
func DeleteNoteByID(conn *pgx.Conn, noteID int) error {
    _, err := conn.Exec(context.Background(), "DELETE FROM Notes WHERE ID = $1", noteID)
    if err != nil {
        return err
    }

    return nil
}
