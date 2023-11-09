package main

import (
    "testing"
	"time"
    "fmt"
	"strings"
	"reflect"
    "sync"
    "github.com/stretchr/testify/assert"
)

// Define all the structs used in the tests
type TestDisplayNote struct {
    NoteID                  int
    NoteTitle               string
    CreationDate            time.Time
    Delegation              string
    CompletionDate          time.Time
    Status                  string
    NoteContent             string
}

type TestUser struct {
    UserID   int    
    UserName string 
}

type MockDB struct {
    notes map[int]TestDisplayNote // Simulate storage for notes
    ErrorOnCreateNote bool // Simulate an error when creating a note
    mu    sync.Mutex
}

type MockDatabaseDelegations struct {
    Delegations []string
    Error       error
}

type MockDatabaseShareList struct {
    Users []TestUser
    Error error
}

type TestCustomSharingList struct {
    ListID   int    
    ListName string 
}


// Constructor function for the mock database
func NewMockDB() *MockDB {
    return &MockDB{
        notes: make(map[int]TestDisplayNote),
    }
}

func (db *MockDB) CreateNote(note TestDisplayNote) error {
     // Protect the critical section with a mutex
     db.mu.Lock()
     defer db.mu.Unlock()
    
    // Check for an empty NoteTitle
    if note.NoteTitle == "" {
        return fmt.Errorf("NoteTitle cannot be empty")
    }

    // Check if an error should be simulated
    if db.ErrorOnCreateNote {
        return fmt.Errorf("Simulated database error")
    }


    // Check for duplicate NoteTitle or other conditions as needed
    for _, existingNote := range db.notes {
        if existingNote.NoteTitle == note.NoteTitle {
            return fmt.Errorf("Note with the same title already exists")
        }
    }
    // Store the note without modifying the NoteID
    db.notes[note.NoteID] = note

    return nil
}


//getNoteByID function to retrieve a note by ID
func (db *MockDB) GetNoteByID(noteID int) (TestDisplayNote, error) {
    // Simulate retrieving a note by ID
    note, ok := db.notes[noteID]
    if !ok {
        return TestDisplayNote{}, fmt.Errorf("Note not found")
    }
    return note, nil
}

//App struct to encapsulate the database dependency
type TestApp struct {
    DB *MockDB
}



//--------create notes test cases-----------------
//test cases for successfully creating a note.
func TestCreateNote_Success(t *testing.T) {
    
	db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // test data for a valid note
    validNote := TestDisplayNote{
        NoteTitle:     "Test Note",
        Delegation:    "John",
        Status:        "Pending",
        NoteContent:   "This is a test note.",
    }

    //createNote function with valid data
    err := app.DB.CreateNote(validNote)

    //no error occurred (note creation was successful)
    assert.Nil(t, err, "Expected no error, but got an error")
}

func TestCreateNote_Duplicate(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Create an initial note
    initialNote := TestDisplayNote{
        NoteTitle:     "Duplicate Note",
        Delegation:    "Alice",
        Status:        "Pending",
        NoteContent:   "This is the initial note.",
    }

    // Create the initial note
    err := app.DB.CreateNote(initialNote)
    assert.Nil(t, err, "Expected no error when creating the initial note")

    // Attempt to create a duplicate note with the same NoteTitle
    duplicateNote := TestDisplayNote{
        NoteTitle:     "Duplicate Note",
        Delegation:    "Bob",
        Status:        "Pending",
        NoteContent:   "This is a duplicate note.",
    }

    // Try to create the duplicate note
    err = app.DB.CreateNote(duplicateNote)

    // Ensure that an error occurred, indicating a duplicate note creation attempt
    assert.NotNil(t, err, "Expected an error for duplicate note creation")
}

func TestCreateNote_InvalidNoteTitle(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

     // Attempt to create a note with an empty NoteTitle
     invalidNote := TestDisplayNote{
        NoteTitle: "", // Empty NoteTitle
    }

    err := app.DB.CreateNote(invalidNote)

    
    // Ensure that an error occurred
    if err == nil {
        t.Error("Expected an error for note creation with an empty NoteTitle, but no error occurred")
    }
}

/*func TestCreateNote_MissingDelegation(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Create a note without specifying a delegation
    noteWithoutDelegation := TestDisplayNote{
        NoteTitle:     "Note Without Delegation",
        Status:        "Pending",
        NoteContent:   "This is a note without delegation.",
    }

    err := app.DB.CreateNote(noteWithoutDelegation)

    // Ensure that no error occurred (note creation was successful)
    assert.Nil(t, err, "Expected no error, but got an error")

    // Optionally, you can retrieve the created note and verify its attributes
    createdNote, err := app.DB.GetNoteByID(1) // Use the appropriate NoteID, e.g., 1
    assert.Nil(t, err, "Expected no error when retrieving the created note")

    // Verify that the created note has the expected attributes
    assert.Equal(t, noteWithoutDelegation.NoteTitle, createdNote.NoteTitle)
    assert.Equal(t, noteWithoutDelegation.Status, createdNote.Status)
    assert.Equal(t, noteWithoutDelegation.NoteContent, createdNote.NoteContent)
    // Ensure that the delegation is empty, as it was not specified
    assert.Empty(t, createdNote.Delegation)
}*/

/*
func TestCreateNote_InvalidStatus(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Attempt to create a note with an invalid status
    invalidStatusNote := TestDisplayNote{
        NoteTitle:     "Invalid Status Note",
        Delegation:    "Alice",
        Status:        "InvalidStatus", // This is an unsupported status
        NoteContent:   "This is a note with an invalid status.",
    }

    err := app.DB.CreateNote(invalidStatusNote)

    // Ensure that an error occurred, indicating an invalid or unsupported status
    assert.NotNil(t, err, "Expected an error for note creation with an invalid status")
}
*/

func TestCreateNote_EmptyNoteContent(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Create a note with empty NoteContent
    emptyContentNote := TestDisplayNote{
        NoteTitle:     "Empty Content Note",
        Delegation:    "Alice",
        Status:        "Pending",
        NoteContent:   "", // Empty NoteContent
    }

    err := app.DB.CreateNote(emptyContentNote)

    // Ensure that no error occurred when creating a note with empty NoteContent
    assert.Nil(t, err, "Expected no error for note creation with empty NoteContent")
}

func TestCreateNote_LongTitleAndContent(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Create a note with a long NoteTitle and NoteContent
    longTitle := "A very long title that exceeds the normal length of a note title and is used to test the function's ability to handle large input data"
    longContent := "A very long note content that exceeds the normal length of a note's content and is used to test the function's ability to handle large input data."

    longNote := TestDisplayNote{
        NoteTitle:     longTitle,
        Delegation:    "Alice",
        Status:        "Pending",
        NoteContent:   longContent,
    }

    err := app.DB.CreateNote(longNote)

    // Ensure that no error occurred when creating a note with long Title and Content
    assert.Nil(t, err, "Expected no error for note creation with long Title and Content")
}

// Test case to simulate a database error
// Test case to simulate a database error
func TestCreateNote_DatabaseError(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Simulate a database error
    db.ErrorOnCreateNote = true

    // Attempt to create a note
    note := TestDisplayNote{
        NoteTitle:     "Test Note",
        Delegation:    "Alice",
        Status:        "Pending",
        NoteContent:   "This is a test note.",
    }

    err := app.DB.CreateNote(note)

    // Ensure that an error occurred as expected
    assert.NotNil(t, err, "Expected an error when the database returns an error")
}


func TestConcurrentNoteCreation(t *testing.T) {
    numGoroutines := 5

    for i := 0; i < numGoroutines; i++ {
        db := NewMockDB() // Create a new database instance for each goroutine

        go func(i int) {
            app := TestApp{
                DB: db,
            }

            note := TestDisplayNote{
                NoteTitle:   fmt.Sprintf("Concurrent Note %d", i),
                Delegation:  "User",
                Status:      "Pending",
                NoteContent: "This is a concurrent note.",
            }

            err := app.DB.CreateNote(note)

            if err != nil {
                t.Errorf("Goroutine %d: Error creating note: %v", i, err)
            }
        }(i)
    }
}











//---------------delete notes test cases-----------------
func (db *MockDB) DeleteNote(noteID int) error {
    // Check if the note exists, and if it does, delete it
    _, ok := db.notes[noteID]
    if !ok {
        return fmt.Errorf("Note not found")
    }
    delete(db.notes, noteID)
    return nil
}

func TestDeleteNote_Success(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Create a test note to delete
    testNote := TestDisplayNote{
        NoteID: 1,
        NoteTitle: "Test Note", // Provide a valid NoteTitle
        // Add other necessary fields
    }

    db.CreateNote(testNote)
    err := app.DB.DeleteNote(testNote.NoteID)

    assert.Nil(t, err, "Expected no error, but got an error")

    _, err = app.DB.GetNoteByID(testNote.NoteID)
    assert.NotNil(t, err, "Expected an error for a deleted note, but no error occurred")
}

func TestDeleteNote_NonExistentNote(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    nonExistentNoteID := 42 // assuming note ID 42 does not exist in the mock DB

    //call deleteNote function with the non-existent note ID
    err := app.DB.DeleteNote(nonExistentNoteID)

    //assert that an error occurred (note not found)
    assert.NotNil(t, err, "Expected an error for attempting to delete a non-existent note, but got no error")

    //optionally, you can check the error message to ensure it's what you expect
    expectedErrorMessage := "Note not found"
    assert.Equal(t, expectedErrorMessage, err.Error(), "Error message does not match the expected message")
}

//------------------update notes test cases----------------
func (db *MockDB) UpdateNote(noteID int, updatedNote TestDisplayNote) error {
    // Check if the note exists
    existingNote, ok := db.notes[noteID]
    if !ok {
        return fmt.Errorf("Note not found")
    }

    existingNote.NoteTitle = updatedNote.NoteTitle
    existingNote.Delegation = updatedNote.Delegation
    existingNote.Status = updatedNote.Status
    existingNote.NoteContent = updatedNote.NoteContent

    // Store the updated note back in the map
    db.notes[noteID] = existingNote

    return nil
}

func (db *MockDB) TestUpdateNote(noteID int, updatedNote TestDisplayNote) error {
    // Check if the note exists
    existingNote, ok := db.notes[noteID]
    if !ok {
        return fmt.Errorf("Note not found")
    }

    // Check for invalid data, for example, empty fields
    if updatedNote.NoteTitle == "" || updatedNote.Delegation == "" {
        return fmt.Errorf("Invalid data: Note title and delegation cannot be empty")
    }

    
    existingNote.NoteTitle = updatedNote.NoteTitle
    existingNote.Delegation = updatedNote.Delegation
    existingNote.Status = updatedNote.Status
    existingNote.NoteContent = updatedNote.NoteContent

    // Store the updated note back in the map
    db.notes[noteID] = existingNote

    return nil
}

func TestUpdateNote_Success(t *testing.T) {
    // Initialize the mock database
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Create a test note to update
    testNote := TestDisplayNote{
        NoteID: 1, // Existing note ID
        NoteTitle: "Original Title",
        Delegation: "John",
        Status: "Pending",
        NoteContent: "Original content",
    }

    db.CreateNote(testNote)

    // Define updated data
    updatedNoteData := TestDisplayNote{
        NoteTitle: "Updated Title",
        Delegation: "Alice",
        Status: "Completed",
        NoteContent: "Updated content",
    }

    // Call the updateNote function with the test note ID and updated data
    err := app.DB.UpdateNote(testNote.NoteID, updatedNoteData)

    // Assert that no error occurred (note update was successful)
    assert.Nil(t, err, "Expected no error, but got an error")

    // Verify that the note was actually updated
    updatedNote, err := app.DB.GetNoteByID(testNote.NoteID)
    assert.Nil(t, err, "Expected no error while fetching the updated note")
    assert.Equal(t, "Updated Title", updatedNote.NoteTitle)
    assert.Equal(t, "Alice", updatedNote.Delegation)
    assert.Equal(t, "Completed", updatedNote.Status)
    assert.Equal(t, "Updated content", updatedNote.NoteContent)
}

func TestUpdateNote_InvalidData(t *testing.T) {
   
    db := NewMockDB()

    // Create a test note to update
    testNote := TestDisplayNote{
        NoteID: 1, // Existing note ID
        NoteTitle: "Original Title",
        Delegation: "John",
        Status: "Pending",
        NoteContent: "Original content",
    }

    db.CreateNote(testNote)

    // Define invalid updated data with empty values
    updatedNoteData := TestDisplayNote{
        NoteTitle: "", // Invalid: Empty title
        Delegation: "", // Invalid: Empty delegation
        Status: "Completed",
        NoteContent: "Updated content",
    }

    // Initialize the app with the mock DB
    app := TestApp{
        DB: db,
    }

    // Call the TestUpdateNote function with the test note ID and invalid data
    err := app.DB.TestUpdateNote(testNote.NoteID, updatedNoteData)
    assert.NotNil(t, err, "Expected an error for updating with invalid data, but got no error")

    expectedErrorMessage := "Invalid data: Note title and delegation cannot be empty"
    assert.Equal(t, expectedErrorMessage, err.Error(), "Error message does not match the expected message")
}

//----------------list notes test cases----------
func (db *MockDB) ListNotes() ([]TestDisplayNote, error) {
    notes := make([]TestDisplayNote, 0)
    for _, note := range db.notes {
        notes = append(notes, note)
    }
    return notes, nil
}

func TestUpdateNote_NonExistentNote(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    // Attempt to update a non-existent note
    nonExistentNoteID := 32 // Assuming note ID 42 does not exist in the mock DB

    // Define updated data
    updatedNoteData := TestDisplayNote{
        NoteTitle: "Updated Title",
        Delegation: "Alice",
        Status: "Completed",
        NoteContent: "Updated content",
    }

    // Call the updateNote function with the non-existent note ID and updated data
    err := app.DB.UpdateNote(nonExistentNoteID, updatedNoteData)
    assert.NotNil(t, err, "Expected an error for attempting to update a non-existent note, but got no error")
}

func TestListNotes_Success(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    note1 := TestDisplayNote{NoteID: 1, NoteTitle: "Note 1", }
    note2 := TestDisplayNote{NoteID: 2, NoteTitle: "Note 2", }
    db.CreateNote(note1)
    db.CreateNote(note2)

    notes, err := app.DB.ListNotes()

    assert.Nil(t, err, "Expected no error, but got an error")
    assert.Contains(t, notes, note1, "List of notes should contain note1")
    assert.Contains(t, notes, note2, "List of notes should contain note2")
}

func TestListNotes_WithNotes(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    testNotes := []TestDisplayNote{
        {NoteID: 1, NoteTitle: "Title 1", Delegation: "John", Status: "Pending", NoteContent: "Content 1"},
        {NoteID: 2, NoteTitle: "Title 2", Delegation: "Alice", Status: "Completed", NoteContent: "Content 2"},
    }

    for _, note := range testNotes {
        db.CreateNote(note)
    }

    notes, err := app.DB.ListNotes()
    assert.Nil(t, err, "Expected no error, but got an error")
    assert.ElementsMatch(t, testNotes, notes, "Listed notes do not match the expected notes")
}

func TestListNotes_EmptyList(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    notes, err := app.DB.ListNotes()
    assert.Nil(t, err, "Expected no error, but got an error")
    assert.Empty(t, notes, "Expected an empty list of notes, but got non-empty list")
}


//-----------------search notes test cases-----------------
func (db *MockDB) SearchNotes(query string) ([]TestDisplayNote, error) {
    matchedNotes := make([]TestDisplayNote, 0)
    for _, note := range db.notes {
        // Perform a case-insensitive search based on NoteTitle and NoteContent
        if strings.Contains(strings.ToLower(note.NoteTitle), strings.ToLower(query)) ||
            strings.Contains(strings.ToLower(note.NoteContent), strings.ToLower(query)) {
            matchedNotes = append(matchedNotes, note)
        }
    }
    return matchedNotes, nil
}


func TestSearchNotes_Success(t *testing.T) {
    
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    note1 := TestDisplayNote{NoteID: 1, NoteTitle: "Important Meeting", NoteContent: "Discuss project deadlines", /* other fields */}
    note2 := TestDisplayNote{NoteID: 2, NoteTitle: "Grocery List", NoteContent: "Buy milk and eggs", /* other fields */}
    db.CreateNote(note1)
    db.CreateNote(note2)

    
    searchQuery := "Meeting"
    matchedNotes, err := app.DB.SearchNotes(searchQuery)
    assert.Nil(t, err, "Expected no error, but got an error")


    assert.Equal(t, len(matchedNotes), 1, "Expected 1 matching note")
    assert.Contains(t, matchedNotes, note1, "Expected note1 to match the search query")
}

func TestSearchNotes_NoMatches(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    note1 := TestDisplayNote{NoteID: 1, NoteTitle: "Important Meeting", NoteContent: "Discuss project deadlines", /* other fields */}
    note2 := TestDisplayNote{NoteID: 2, NoteTitle: "Grocery List", NoteContent: "Buy milk and eggs", /* other fields */}
    db.CreateNote(note1)
    db.CreateNote(note2)

    searchQuery := "Vacation"
    matchedNotes, err := app.DB.SearchNotes(searchQuery)

    assert.Nil(t, err, "Expected no error, but got an error")

    assert.Empty(t, matchedNotes, "Expected no matching notes for the search query")
}

func TestSearchNotes_EmptyQuery(t *testing.T) {
    db := NewMockDB()

    app := TestApp{
        DB: db,
    }

    
    note1 := TestDisplayNote{NoteID: 1, NoteTitle: "Important Meeting", NoteContent: "Discuss project deadlines", /* other fields */}
    note2 := TestDisplayNote{NoteID: 2, NoteTitle: "Grocery List", NoteContent: "Buy milk and eggs", /* other fields */}
    db.CreateNote(note1)
    db.CreateNote(note2)

    // Perform a search with an empty query
    searchQuery := ""
    matchedNotes, err := app.DB.SearchNotes(searchQuery)

    
    assert.Nil(t, err, "Expected no error, but got an error")

    expectedNotes := []TestDisplayNote{note1, note2}
    assert.ElementsMatch(t, expectedNotes, matchedNotes, "Expected all notes to be returned for an empty query")
}


// ------------ delegations test case --------------
func getDelegationsHandler(db MockDatabaseDelegations) ([]string, error) {
    // Simulate retrieving delegations from the mock database.
    return db.Delegations, db.Error
}

func TestGetDelegationsHandler(t *testing.T) {
    
    mockDB := MockDatabaseDelegations{
        Delegations: []string{"Delegation1", "Delegation2"},
        Error:       nil,
    }

    delegations, err := getDelegationsHandler(mockDB)

    // Check the result
    if err != nil {
        t.Errorf("Expected no error, but got an error: %v", err)
    }

    expectedDelegations := []string{"Delegation1", "Delegation2"}
    if !reflect.DeepEqual(delegations, expectedDelegations) {
        t.Errorf("Expected delegations %v, but got %v", expectedDelegations, delegations)
    }
}

func TestGetDelegationsHandler_Success(t *testing.T) {
    
    mockDB := MockDatabaseDelegations{
        Delegations: []string{"Delegation1", "Delegation2"},
        Error:       nil,
    }

    delegations, err := getDelegationsHandler(mockDB)

    // Assert that no error occurred
    assert.Nil(t, err, "Expected no error, but got an error: %v", err)

    // Define the expected list of delegations
    expectedDelegations := []string{"Delegation1", "Delegation2"}

    assert.ElementsMatch(t, expectedDelegations, delegations, "Delegations do not match the expected list")
}

func TestGetDelegationsHandler_EmptyList(t *testing.T) {
   
    mockDB := MockDatabaseDelegations{
        Delegations: []string{},
        Error:       nil,
    }

    delegations, err := getDelegationsHandler(mockDB)

    assert.Nil(t, err, "Expected no error, but got an error: %v", err)
    assert.Empty(t, delegations, "Expected an empty list of delegations")
}

func TestGetDelegationsHandler_Error(t *testing.T) {
    mockDB := MockDatabaseDelegations{
        Delegations: nil,
        Error:       fmt.Errorf("Database error"),
    }

    delegations, err := getDelegationsHandler(mockDB)

    assert.NotNil(t, err, "Expected an error, but got no error")
    assert.Nil(t, delegations, "Expected delegations to be nil due to the error")
}

// ------------ share list test case ---------------
func getShareListHandler(db MockDatabaseShareList) ([]TestUser, error) {

	// Simulate retrieving users from the mock database.
	if db.Error != nil {
        return nil, db.Error
    }

    users := make([]TestUser, len(db.Users))
    for i, user := range db.Users {
        owner := TestUser{
            UserID:   user.UserID,
            UserName: user.UserName, // Corrected field name
        }
        users[i] = owner
    }
     return users, nil
}

func TestGetShareListHandler(t *testing.T) {
    
    mockDB := MockDatabaseShareList{
        Users: []TestUser{
            {UserID: 1, UserName: "User1"}, 
            {UserID: 2, UserName: "User2"}, 
        },
        Error: nil,
    }

    owners, err := getShareListHandler(mockDB)

    if err != nil {
        t.Errorf("Expected no error, but got an error: %v", err)
    }

    // Check the returned owner data (you may need to compare with expected data)
    expectedOwners := []TestUser{
        {UserID: 1, UserName: "User1"}, 
        {UserID: 2, UserName: "User2"}, 
    }
    if !reflect.DeepEqual(owners, expectedOwners) {
        t.Errorf("Expected owners %v, but got %v", expectedOwners, owners)
    }
}

func TestGetShareListHandler_Success(t *testing.T) {

    mockDB := MockDatabaseShareList{
        Users: []TestUser{
            {UserID: 1, UserName: "User1"},
            {UserID: 2, UserName: "User2"},
        },
        Error: nil,
    }
    owners, err := getShareListHandler(mockDB)

    assert.Nil(t, err, "Expected no error, but got an error: %v", err)

    // Define the expected list of owners
    expectedOwners := []TestUser{
        {UserID: 1, UserName: "User1"},
        {UserID: 2, UserName: "User2"},
    }
    assert.ElementsMatch(t, expectedOwners, owners, "Owners do not match the expected list")
}

func TestGetShareListHandler_EmptyList(t *testing.T) {
    
    mockDB := MockDatabaseShareList{
        Users: []TestUser{},
        Error: nil,
    }

    owners, err := getShareListHandler(mockDB)

    assert.Nil(t, err, "Expected no error, but got an error: %v", err)
    assert.Empty(t, owners, "Expected an empty list of owners")
}

func TestGetShareListHandler_Error(t *testing.T) {

    mockDB := MockDatabaseShareList{
        Users: nil,
        Error: fmt.Errorf("Database error"),
    }

    owners, err := getShareListHandler(mockDB)

    assert.NotNil(t, err, "Expected an error, but got no error")
    assert.Nil(t, owners, "Expected owners to be nil due to the error")
}


//-----------------custom sharing list test cases-----------------
func (db *MockDB) GetCustomSharingLists(userID int) ([]TestCustomSharingList, error) {
	if userID == 1 {
        return []TestCustomSharingList{
            {ListID: 1, ListName: "List1"},
            {ListID: 2, ListName: "List2"},
        }, nil
    } else if userID == 3 {
        return nil, fmt.Errorf("Invalid user ID")
    }
    return []TestCustomSharingList{}, nil
}

func TestGetCustomSharingLists_Success(t *testing.T) {
    
    db := &MockDB{
        notes: make(map[int]TestDisplayNote),
    }

    app := TestApp{
        DB: db,
    }

    lists, err := app.DB.GetCustomSharingLists(1)

   
    assert.Nil(t, err, "Expected no error, but got an error")

    expectedLists := []TestCustomSharingList{
        {ListID: 1, ListName: "List1"},
        {ListID: 2, ListName: "List2"},
    }
    assert.ElementsMatch(t, lists, expectedLists, "Custom sharing lists do not match the expected lists")
}

func TestGetCustomSharingLists_NoLists(t *testing.T) {
    db := &MockDB{}

    app := TestApp{
        DB: db,
    }

    lists, err := app.DB.GetCustomSharingLists(2) // Assuming user 2 has no custom sharing lists

    assert.Nil(t, err, "Expected no error, but got an error: %v", err)
    assert.Empty(t, lists, "Expected an empty list of custom sharing lists")
}

func TestGetCustomSharingLists_InvalidUserID(t *testing.T) {
    
    db := &MockDB{}

    app := TestApp{
        DB: db,
    }

    lists, err := app.DB.GetCustomSharingLists(3) // Assuming user 3 does not exist

    assert.NotNil(t, err, "Expected an error for an invalid user ID, but got no error")
    assert.Nil(t, lists, "Expected custom sharing lists to be nil due to the error")
}