package main

import (
	"encoding/json"
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
	Notes    []DisplayNote
}

type DisplayNote struct {
	NoteID                  int
	NoteTitle               string
	CreationDate            time.Time
	Delegation              string
	CompletionDate          time.Time
	Status                  string
	Username                string
	NoteContent             string
	CreationDateFormatted   string
	CompletionDateFormatted string
}

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("index")
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("list")
	a.isAuthenticated(w, r)

	// get the current username
	sess := session.Get(r)
	log.Printf("Session received")

	user := "[guest]"
	log.Printf("Temp user made")

	// Check if there is a session and retrieve the username if available
	if sess != nil {
		user = sess.CAttr("username").(string)
		log.Printf("User updated")
	}

	if r.Method != "GET" {
		// Handle incorrect HTTP method
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

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
		SQL = `SELECT 
            notes.noteID, 
            notes.note_title, 
            notes.creationdate, 
            notes.delegatedto, 
            notes.completion_date,
            notes.status, 
            users.username, 
            notes.note_content
            FROM "notes"
            JOIN users ON notes.userid = users.userid
            ORDER by notes.noteID;`
	}

	// Execute the SQL query to retrieve notes
	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)
	log.Println("Query Executed")

	// Define a function map for use in the template.
	var funcMap = template.FuncMap{
		"addOne": func(n int) int {
			return n + 1
		},
	}

	// Create a Data structure to pass to the template
	data := Data{}
	data.Username = user
	var note DisplayNote
	// Loop through the rows and scan note information from the database.
	for rows.Next() {
		err := rows.Scan(&note.NoteID, &note.NoteTitle, &note.CreationDate, &note.Delegation, &note.CompletionDate, &note.Status, &note.Username, &note.NoteContent)
		checkInternalServerError(err, w)
		note.CreationDateFormatted = formatDateForMainPage(note.CreationDate)
		note.CompletionDateFormatted = formatDateForMainPage(note.CompletionDate)
		checkInternalServerError(err, w)
		data.Notes = append(data.Notes, note)
	}

	// Load the template and execute it with the data
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

func formatDateForMainPage(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func (a *App) createHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("create")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var note Note
	sess := session.Get(r)
	log.Println(sess)
	note.UserID = sess.CAttr("userid").(int)
	log.Printf("Note: userid")
	log.Println(note)
	note.NoteTitle = r.FormValue("NoteTitle")
	log.Printf("Note: notetitle")
	log.Println(note)
	note.NoteContent = r.FormValue("NoteContent")
	log.Printf("Note: note content")
	log.Println(note)
	note.CreationDate, err = time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
	log.Printf("Note: creatioin date")
	log.Println(note)
	note.DelegatedTo = r.FormValue("delegated")
	log.Printf("Note: delegated to")
	log.Println(note)
	note.CompletionDate, err = time.Parse("2006-01-02T15:04", r.FormValue("CompletionDate"))
	if err != nil {
		log.Printf("Error with Completion Date")
		checkInternalServerError(err, w)
	}
	log.Printf("Note: completeion date")
	log.Println(note)
	note.Status = r.FormValue("status")
	log.Printf("Note: ststaus")
	log.Println(note)

	// Prepare an SQL statement to insert the new note
	stmt, err := a.db.Prepare(`INSERT INTO "notes"(userID, note_title, note_content, creationdate, delegatedto, completion_date, status) VALUES($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		log.Printf("Prepare query error")
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
	log.Println(sess)
	note.UserID = sess.CAttr("userid").(int)
	log.Println(note)
	note.NoteContent = r.FormValue("NoteContent")
	log.Println(note)
	note.DelegatedTo = r.FormValue("delegated")
	log.Println(note)
	note.Status = r.FormValue("status")
	log.Println(note)
	completionDateStr := r.FormValue("completiondate")
	log.Println(completionDateStr)
	CompletionDate, err := time.Parse("2006-01-02T15:04", completionDateStr)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid completion date", http.StatusBadRequest)
		return
	}
	note.CompletionDate = CompletionDate
	log.Println(note)
	note.Status = r.FormValue("status")
	log.Println(note)
	log.Println(r.FormValue("noteIdToUpdate"))
	note.NoteID, err = strconv.Atoi(r.FormValue("noteIdToUpdate"))
	log.Println(note)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invaild note id", http.StatusBadRequest)
		return
	}

	stmt, err := a.db.Prepare(`UPDATE "notes" SET  note_content=$1, status=$2, delegatedto=$3, completion_date=$4 WHERE noteID=$5`)
	checkInternalServerError(err, w)

	res, err := stmt.Exec(note.NoteContent, note.Status, note.DelegatedTo, note.CompletionDate, note.NoteID)
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
	var noteID, _ = strconv.ParseInt(r.FormValue("NoteId"), 10, 64)
	stmt, err := a.db.Prepare(`DELETE FROM "notes" WHERE noteID=$1`)
	checkInternalServerError(err, w)

	res, err := stmt.Exec(noteID)
	checkInternalServerError(err, w)

	_, err = res.RowsAffected()
	checkInternalServerError(err, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *App) getdelegationsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Delegations")
	names := []string{}
	SQL := `SELECT username from "users"`
	rows, err := a.db.Query(SQL)
	//log.Println(rows)
	checkInternalServerError(err, w)
	var name string
	for rows.Next() {
		err := rows.Scan(&name)
		//log.Println(name)
		checkInternalServerError(err, w)
		names = append(names, name)
		//log.Println(names)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}
