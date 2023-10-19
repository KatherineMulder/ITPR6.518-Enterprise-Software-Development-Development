package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/icza/session"
)

// hold user data and notes for rendering templates.
type Data struct {
	Username string
	Notes    []Note
	//Sharing []Note
	
}

// The indexHandler handles the root endpoint.
func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("index")
	a.isAuthenticated(w, r)
	// Redirect to the list page after authentication
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

// The listHandler handles the listing of notes.
func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("list")

	// Check if the user is authenticated
	a.isAuthenticated(w, r) 

	//get the current username
	sess := session.Get(r)
	log.Println(sess)

	user := "[guest]"
	log.Println(user)

	// Check if there is a session and retrieve the username if available
	if sess != nil {
		user = sess.CAttr("username").(string)
	}

	if r.Method != "GET" {
		// Handle incorrect HTTP method
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}


	params := mux.Vars(r)
	log.Println(params)
	sortcol, err := strconv.Atoi(params["srt"])
	log.Println(sortcol)

	_, ok := params["srt"]
	if ok && err != nil {
		// Redirect to the list page if the sort parameter is invalid
		http.Redirect(w, r, "/list", http.StatusFound)
	}

	SQL := ""
	log.Println(SQL)
	switch sortcol {
	case 1:
		SQL = `SELECT * FROM "notes" ORDER by note_title`
	case 2:
		SQL = `SELECT * FROM "notes" ORDER by creation_date`
	case 3:
		SQL = `SELECT * FROM "notes" ORDER by completion_date`
	case 4:
		SQL = `SELECT * FROM "notes" ORDER by status`
	default:
		SQL = `SELECT * FROM "notes" ORDER by noteID`
	}
	log.Println(SQL)
	rows, err := a.db.Query(SQL)
	log.Println(rows)

	//// Check for internal server errors and handle them by writing an error response.
	checkInternalServerError(err, w)

	// Define a function map for use in the template.
	var funcMap = template.FuncMap{
		"addOne": func(n int) int {
			return n + 1
		},
	}
	log.Println(funcMap)

	// Initialize the data structure to hold information for the template.
	data := Data{}
	log.Println(data)

	// Set the username in the data structure.
	data.Username = user
	log.Println(data)

	// Initialize a note variable.
	var note Note
	log.Println(note)

	// Loop through the rows and scan note information from the database.
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.UserID, &note.NoteTitle, &note.NoteContent, &note.CompletionDate, &note.Status)
		checkInternalServerError(err, w)
		log.Println(err)
		note.FormattedDate()
		checkInternalServerError(err, w)

		// Append the scanned note to the data structure.
		data.Notes = append(data.Notes, note)
		log.Println(data)
	}
	// Create a new template, specify the function map, and parse the template file.
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

// The createHandler handles creating a new note.
func (a *App) createHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("create")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var note Note
	sess := session.Get(r)
	note.UserID = sess.CAttr("userID").(int)
	note.NoteTitle = r.FormValue("NoteTitle")
	note.NoteContent = r.FormValue("NoteContent")
	if err != nil {
		log.Fatal(err)
	}
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate")) // // Parse the CompletionDate using the specified format
	if err != nil {
		log.Fatal(err)
	}
	note.Status = r.FormValue("status")

	// Prepare an SQL statement to insert the new note
	stmt, err := a.db.Prepare(`INSERT INTO "notes"(userID, note_title, note_content, creation_date, completetion_date, status) VALUES($1, $2, $3, $4, $5, $6)`)

	if err != nil {
		// Log and handle any errors related to SQL statement preparation
		log.Printf("Error with Query Prepare")
		checkInternalServerError(err, w)
	}
	_, err = stmt.Exec(note.UserID, note.NoteTitle, note.NoteContent, note.CompletionDate, note.Status)
	if err != nil {
		// Log and handle any errors related to SQL statement execution
		log.Printf("Error with Executing Query")
		checkInternalServerError(err, w)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// // The updateHandler handles updating a note. ////
func (a *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("update")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var note Note
	note.NoteID, _ = strconv.Atoi(r.FormValue("NoteID"))
	note.NoteTitle = r.FormValue("NoteTitle")
	note.NoteContent = r.FormValue("NoteContent")
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate")) //// Parse the CompletionDate using the specified format
	if err != nil {
		log.Fatal(err)
	}
	note.Status = r.FormValue("status")

	// Prepare an SQL statement to update the note
	stmt, err := a.db.Prepare(`UPDATE "notes" SET note_title=$1, note_content=$2, completion_date=$3, status=$4 WHERE noteID=$5`)
	checkInternalServerError(err, w)

	// Execute the SQL statement to update the note
	res, err := stmt.Exec(note.NoteTitle, note.NoteContent, note.CompletionDate, note.Status, note.NoteID)
	checkInternalServerError(err, w)

	// Check the number of rows affected by the update
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// //  The deleteHandler handles deleting a note. ////
func (a *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Delete")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var noteID, _ = strconv.ParseInt(r.FormValue("NoteID"), 10, 64) // Parse the noteID from the form data
	stmt, err := a.db.Prepare(`DELETE FROM "notes" WHERE noteID=$1`)
	checkInternalServerError(err, w)

	// Execute the SQL statement to delete the note
	res, err := stmt.Exec(noteID)
	checkInternalServerError(err, w)

	// Check the number of rows affected by the deletion
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
