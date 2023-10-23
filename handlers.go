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

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("index")
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

/*the list handler process:
1.Authentication Check
2.Session Handling
3.HTTP Method Check
4.Data Retrieval
5.Shared Users for Each Note
6.Data Preparation
7.Template Rendering
8.HTTP response.
*/

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("list")
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

	// ======= get all notes from the "notes" table ========== //
	// Determine the sorting index
	params := mux.Vars(r)
	sortcol, err := strconv.Atoi(params["srt"])
	_, ok := params["srt"]
	if ok && err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	SQL := ""
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

	// Execute the SQL query to retrieve notes
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)

	// Define a function map for use in the template.
	var funcMap = template.FuncMap{
		"addOne": func(n int) int {
			return n + 1
		},
	}

	// Create a Data structure to pass to the template
	data := Data{}
	data.Username = user

	// Loop through the rows and scan note information from the database.
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.NoteID, &note.UserID, &note.NoteTitle, &note.NoteContent, &note.CreationDate, &note.DelegatedTo, &note.CompletionDate, &note.Status)
		checkInternalServerError(err, w)
		note.FormattedDate()
		checkInternalServerError(err, w)
		data.Notes = append(data.Notes, note)
	}

	// Load the template and execute it with the data
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

//get all users
//get all shared notes with privileges

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
	note.CreationDate, err = time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
	note.DelegatedTo = r.FormValue("DelegatedTo")
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate"))
	note.Status = r.FormValue("status")

	if err != nil {
		log.Fatal(err)
	}
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate"))
	if err != nil {
		log.Fatal(err)
	}
	note.Status = r.FormValue("status")

	// Prepare an SQL statement to insert the new note
	stmt, err := a.db.Prepare(`INSERT INTO "notes"(userID, note_title, note_content, creation_date, delegated_to, completetion_date, status) VALUES($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		// Log and handle any errors related to SQL statement preparation
		log.Printf("Error with Query Prepare")
		checkInternalServerError(err, w)
	}

	_, err = stmt.Exec(note.UserID, note.NoteTitle, note.NoteContent, note.CreationDate, note.DelegatedTo, note.CompletionDate, note.Status)
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
	sess := session.Get(r)
	note.UserID = sess.CAttr("userID").(int)
	note.NoteTitle = r.FormValue("NoteTitle")
	note.NoteContent = r.FormValue("NoteContent")
	note.CreationDate, err = time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
	note.DelegatedTo = r.FormValue("DelegatedTo")
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate"))
	note.Status = r.FormValue("status")

	if err != nil {
		log.Fatal(err)
	}
	note.Status = r.FormValue("status")

	stmt, err := a.db.Prepare(`UPDATE "notes" SET note_title=$1, note_content=$2, completion_date=$3, status=$4 WHERE noteID=$5`)
	checkInternalServerError(err, w)

	res, err := stmt.Exec(note.NoteTitle, note.NoteContent, note.CompletionDate, note.Status, note.NoteID)
	checkInternalServerError(err, w)

	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Delete")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var noteID, _ = strconv.ParseInt(r.FormValue("NoteID"), 10, 64)
	stmt, err := a.db.Prepare(`DELETE FROM "notes" WHERE noteID=$1`)
	checkInternalServerError(err, w)

	res, err := stmt.Exec(noteID)
	checkInternalServerError(err, w)

	_, err = res.RowsAffected()
	checkInternalServerError(err, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
