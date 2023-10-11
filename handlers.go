package main

import (
	"html/template"
	"net/http"
	"strconv"

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

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	a.isAuthenticated(w, r)
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}
