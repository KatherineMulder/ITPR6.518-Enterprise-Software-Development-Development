package databasesetup

import (
    "testing"
	"context"
    "time"
)

func TestCRUDOperations(t *testing.T) {
    conn, err := DatabaseSetup()
    if err != nil {
        t.Fatal("Error setting up the database:", err)
    }
    defer conn.Close(context.Background())

    t.Run("CreateNote", func(t *testing.T) {
        err := CreateNote(conn, 1, "Sample Note Title", "Sample Note Content", time.Now(), time.Now(), "Active", 1)
        if err != nil {
            t.Fatal("Failed to create a new note:", err)
        }
    })

    t.Run("GetNoteByID", func(t *testing.T) {
        note, err := GetNoteByID(conn, 1)
        if err != nil {
            t.Fatal("Failed to retrieve a note by ID:", err)
        }

        // Add assertions to check the retrieved note's properties
        if note.NoteName != "Sample Note Title" {
            t.Errorf("Expected title: Sample Note Title, got: %s", note.NoteName)
        }

        // Add similar assertions for other properties
    })

    t.Run("UpdateNote", func(t *testing.T) {
        err := UpdateNote(conn, 1, "Updated Title", "Updated Content", time.Now(), "Inactive", 2)
        if err != nil {
            t.Fatal("Failed to update the note:", err)
        }
    })

    t.Run("DeleteNoteByID", func(t *testing.T) {
        err := DeleteNoteByID(conn, 1)
        if err != nil {
            t.Fatal("Failed to delete the note:", err)
        }
    })
}
