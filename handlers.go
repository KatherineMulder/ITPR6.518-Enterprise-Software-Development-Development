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

type Data struct {
	Username string
	Notes    []Note
}

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)

	sess := session.Get(r)
	user := "[guest]"

	if sess != nil {
		user = sess.CAttr("username").(string)
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	params := mux.Vars(r)
	sortcol, err := strconv.Atoi(params["srt"])

	_, ok := params["srt"]
	if ok && err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
	}

	SQL := ""

	switch sortcol {
	case 1:
		SQL = "SELECT * FROM Notes ORDER by userID"
	case 2:
		SQL = "SELECT * FROM Notes ORDER by noteTitle"
	case 3:
		SQL = "SELECT * FROM Notes ORDER by creationDate"
	case 4:
		SQL = "SELECT * FROM Notes ORDER by completionDate"
	case 5:
		SQL = "SELECT * FROM Notes ORDER by status"
	default:
		SQL = "SELECT * FROM Notes ORDER by noteID"
	}

	rows, err := a.db.Query(SQL)
	checkInternalServerError(err, w)
	var funcMap = template.FuncMap{
		"multiplication": func(n int, f int) int {
			return n * f
		},
		"addOne": func(n int) int {
			return n + 1
		},
	}

	data := Data{}
	data.Username = user

	var note Note
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.UserID, &note.NoteTitle, &note.NoteContent, &note.CreationDate, &note.CompletionDate, &note.Status)
		checkInternalServerError(err, w)
		data.Notes = append(data.Notes, note)
	}
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("templates/list.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	checkInternalServerError(err, w)
}

func (a *App) createHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var note Note
	sess := session.Get(r)
	note.UserID = sess.CAttr("userID").(int)
	note.NoteTitle = r.FormValue("NoteTitle")
	note.NoteContent = r.FormValue("NoteContent")
	note.CreationDate, err = time.Parse("2006-01-02 15:04", r.FormValue("createDate"))
	if err != nil {
		log.Fatal(err)
	}
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate"))
	if err != nil {
		log.Fatal(err)
	}
	note.Status = r.FormValue("status")

	stmt, err := a.db.Prepare("INSERT INTO Notes(userID, noteTitle, noteContent, creationDate, completetionDate, status) VALUES($1, $2, $3, $4, $5, $6)")

	if err != nil {
		log.Printf("Error with Query Prepare")
		checkInternalServerError(err, w)
	}
	_, err = stmt.Exec(note.UserID, note.NoteTitle, note.NoteContent, note.CreationDate, note.CompletionDate, note.Status)
	if err != nil {
		log.Printf("Error with Executing Query")
		checkInternalServerError(err, w)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	var note Note
	note.NoteID, _ = strconv.Atoi(r.FormValue("NoteID"))
	note.NoteTitle = r.FormValue("NoteTitle")
	note.NoteContent = r.FormValue("NoteContent")
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate"))
	if err != nil {
		log.Fatal(err)
	}
	note.Status = r.FormValue("status")

	stmt, err := a.db.Prepare("UPDATE Notes SET noteTitle=$1, noteContent=$2, completionDate=$3, status=$4 WHERE noteID=$5")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(note.NoteTitle, note.NoteContent, note.CompletionDate, note.Status, note.NoteID)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
	var noteID, _ = strconv.ParseInt(r.FormValue("NoteID"), 10, 64)
	stmt, err := a.db.Prepare("DELETE FROM Notes WHERE noteID=$1")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(noteID)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}
