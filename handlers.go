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
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("list")
	a.isAuthenticated(w, r)

	// get the current username
	sess := session.Get(r)
	log.Printf("Session received")

	// Set the default username to guest
	user := "[guest]"
	userid := 0
	log.Printf("Temp user made")

	// Check if there is a session and retrieve the username if available
	if sess != nil {
		user = sess.CAttr("username").(string)
		userid = sess.CAttr("userid").(int)
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

	// Redirect to the default sorting index if the sorting index is invalid
	_, ok := params["srt"]
	if ok && err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	// Define the SQL query to retrieve notes
	SQL := ""
	switch sortcol {
	case 1:
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
		JOIN sharing ON notes.noteid = sharing.noteid
		WHERE notes.userid = $1
		OR sharing.userid = $1 
		ORDER by username`
	case 2:
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
		JOIN sharing ON notes.noteid = sharing.noteid
		WHERE notes.userid = $1
		OR sharing.userid = $1 
		ORDER by note_title`
	case 3:
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
		JOIN sharing ON notes.noteid = sharing.noteid
		WHERE notes.userid = $1
		OR sharing.userid = $1
		ORDER by creationdate`
	case 4:
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
		JOIN sharing ON notes.noteid = sharing.noteid
		WHERE notes.userid = $1
		OR sharing.userid = $1
		ORDER by delegatedto`
	case 5:
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
		JOIN sharing ON notes.noteid = sharing.noteid
		WHERE notes.userid = $1
		OR sharing.userid = $1 
		ORDER by completion_date`
	case 6:
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
		JOIN sharing ON notes.noteid = sharing.noteid
		WHERE notes.userid = $1
		OR sharing.userid = $1 ORDER by status`
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
			JOIN sharing ON notes.noteid = sharing.noteid
			WHERE notes.userid = $1
			OR sharing.userid = $1
            ORDER by notes.noteID;`
	}

	// Execute the SQL query to retrieve notes
	rows, err := a.db.Query(SQL, userid)
	log.Println(rows)
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Template loaded")

	// Execute the template
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

func formatDateForMainPage(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func (a *App) searchNotesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("search notes")
	a.isAuthenticated(w, r)
	sess := session.Get(r)
	user := sess.CAttr("username").(string)
	userid := sess.CAttr("userid").(int)

	query := r.FormValue("searchfield")

	rows, err := a.db.Query(`SELECT 
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
	JOIN sharing ON notes.noteid = sharing.noteid
	WHERE (notes.userid = $1
	OR sharing.userid = $1) AND (
	notes.noteID::TEXT ILIKE '%' || UPPER($2) || '%'
   	OR UPPER(notes.note_title) ILIKE '%' || UPPER($2) || '%'
   	OR DATE_TRUNC('minute', notes.creationdate)::TEXT ILIKE '%' || UPPER($2) || '%'
   	OR UPPER(notes.delegatedto) ILIKE '%' || UPPER($2) || '%'
   	OR DATE_TRUNC('minute', notes.completion_date)::TEXT ILIKE '%' || UPPER($2) || '%'
   	OR UPPER(notes.status::TEXT) ILIKE '%' || UPPER($2) || '%'
   	OR UPPER(users.username) ILIKE '%' || UPPER($2) || '%')
	ORDER by notes.noteID;`, userid, query)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 400)
	}

	var funcMap = template.FuncMap{
		"addOne": func(n int) int {
			return n + 1
		},
	}

	data := Data{}
	data.Username = user
	var note DisplayNote
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Template loaded")

	// Execute the template
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

func (a *App) createHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("create")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	// Retrieve the note information from the form
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

	// Set the status to "Not Started" by default
	note.Status = r.FormValue("status")
	log.Printf("Note: ststaus")
	log.Println(note)

	// Prepare an SQL statement to insert the new note
	stmt, err := a.db.Prepare(`INSERT INTO "notes"(userID, note_title, note_content, creationdate, delegatedto, completion_date, status) VALUES($1, $2, $3, $4, $5, $6, $7)`)

	// Check for errors related to the SQL statement
	if err != nil {
		log.Printf("Prepare query error")
		checkInternalServerError(err, w)
	}

	// Execute the SQL statement
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

	// Retrieve the note information from the form
	var note Note
	sess := session.Get(r)
	note.UserID = sess.CAttr("userid").(int)
	note.NoteContent = r.FormValue("NoteContent")
	note.DelegatedTo = r.FormValue("delegated")
	note.Status = r.FormValue("status")
	completionDateStr := r.FormValue("completiondate")
	CompletionDate, err := time.Parse("2006-01-02T15:04", completionDateStr)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid completion date", http.StatusBadRequest)
		return
	}
	formattedDate := CompletionDate.Format("2006-01-02 15:04:05")
	note.Status = r.FormValue("status")
	log.Println(r.FormValue("noteIdToUpdate"))
	note.NoteID, err = strconv.Atoi(r.FormValue("noteIdToUpdate"))
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invaild note id", http.StatusBadRequest)
		return
	}

	stmt, err := a.db.Prepare(`UPDATE "notes" SET  note_content=$1, status=$2, delegatedto=$3, completion_date=$4 WHERE noteID=$5`)
	checkInternalServerError(err, w)

	// Execute the SQL statement
	res, err := stmt.Exec(note.NoteContent, note.Status, note.DelegatedTo, formattedDate, note.NoteID)
	log.Println(res)
	log.Println(err)
	checkInternalServerError(err, w)

	// Check the number of rows affected by the update
	_, err = res.RowsAffected()
	log.Println(err)
	checkInternalServerError(err, w)
}

func (a *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Delete")
	a.isAuthenticated(w, r)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	// Retrieve the noteID from the form
	var noteID, _ = strconv.ParseInt(r.FormValue("NoteId"), 10, 64)
	stmt, err := a.db.Prepare(`DELETE FROM "notes" WHERE noteID=$1`)
	checkInternalServerError(err, w)

	res, err := stmt.Exec(noteID)
	checkInternalServerError(err, w)

	_, err = res.RowsAffected()
	checkInternalServerError(err, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *App) shareNoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Sharing note")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	noteid, err := strconv.Atoi(r.FormValue("noteIdToUpdate-Share"))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(noteid)
	selectedusers := r.Form["user"]
	log.Println(selectedusers)

	for _, userIDStr := range selectedusers {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Printf("Error converting user ID '%s' to an integer: %v", userIDStr, err)
			continue
		}

		SQL, err := a.db.Prepare(`INSERT INTO sharing(noteid, userid, setup_date) Values($1,$2,$3)`)
		if err != nil {
			log.Println(err)
			return
		}
		currenttime := time.Now()

		res, err := SQL.Exec(noteid, userID, currenttime)
		if err != nil {
			log.Println(err)
		}

		res.RowsAffected()
	}

	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

func (a *App) getdelegationsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Delegations")

	// Retrieve the noteID from the form
	names := []string{}
	SQL := `SELECT username from "users"`
	rows, err := a.db.Query(SQL)

	//log.Println(rows)
	checkInternalServerError(err, w)

	// Loop through the rows and scan note information from the database.
	var name string
	for rows.Next() {
		err := rows.Scan(&name)
		//log.Println(name)
		checkInternalServerError(err, w)
		names = append(names, name)
		//log.Println(names)
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func (a *App) getShareListHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Share list")

	// Define a struct to hold the data you want to return
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Initialize a slice to store the user data
	users := []User{}

	// Execute the SQL query to fetch user data
	SQL := `SELECT userid, username from users`
	rows, err := a.db.Query(SQL)
	if err != nil {
		log.Printf("Error querying the database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var user User
	// Loop through the rows and scan user information from the database.
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Printf("Error scanning database row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the user data into JSON and send it as the response
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
