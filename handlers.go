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

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("index")
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("list")
	a.isAuthenticated(w, r)

	sess := session.Get(r)
	log.Println(sess)
	user := "[guest]"
	log.Println(user)

	if sess != nil {
		user = sess.CAttr("username").(string)
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	params := mux.Vars(r)
	log.Println(params)
	sortcol, err := strconv.Atoi(params["srt"])
	log.Println(sortcol)

	_, ok := params["srt"]
	if ok && err != nil {
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

	checkInternalServerError(err, w)
	var funcMap = template.FuncMap{
		"addOne": func(n int) int {
			return n + 1
		},
	}
	log.Println(funcMap)

	data := Data{}
	log.Println(data)
	data.Username = user
	log.Println(data)

	var note Note
	log.Println(note)
	for rows.Next() {
		err = rows.Scan(&note.UserID, &note.NoteTitle, &note.NoteContent, &note.CompletionDate, &note.Status)
		log.Println(err)
		checkInternalServerError(err, w)
		data.Notes = append(data.Notes, note)
		log.Println(data)
	}
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	log.Println(t)
	checkInternalServerError(err, w)
	err = t.Execute(w, data)
	log.Println(t)
	checkInternalServerError(err, w)
}

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
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate"))
	if err != nil {
		log.Fatal(err)
	}
	note.Status = r.FormValue("status")

	stmt, err := a.db.Prepare(`INSERT INTO "notes"(userID, note_title, note_content, creation_date, completetion_date, status) VALUES($1, $2, $3, $4, $5, $6)`)

	if err != nil {
		log.Printf("Error with Query Prepare")
		checkInternalServerError(err, w)
	}
	_, err = stmt.Exec(note.UserID, note.NoteTitle, note.NoteContent, note.CompletionDate, note.Status)
	if err != nil {
		log.Printf("Error with Executing Query")
		checkInternalServerError(err, w)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

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
	note.CompletionDate, err = time.Parse("2006-01-02 15:04", r.FormValue("CompletionDate"))
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
