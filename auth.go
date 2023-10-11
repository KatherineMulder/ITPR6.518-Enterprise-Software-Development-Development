package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/icza/session"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/register/html")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	var user User
	err := a.db.QueryRow("SELECT Username, Password, Role FROM Users WHERE Username=$1", username).Scan(&user.Username, &user.Password, &user.Role)
	switch {
	case err == sql.ErrNoRows:
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)
		_, err = a.db.Exec("INSERT INTO Users(username, password, role) VALUES($1, $2, $3)", username, hashedpassword, role)
		checkInternalServerError(err, w)
	case err != nil:
		http.Error(w, "loi: "+err.Error(), http.StatusBadRequest)
		return
	default:
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method %s", r.Method)
	if r.Method != "POST" {
		http.ServeFile(w, r, "template/login.html")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User
	err := a.db.QueryRow("SELECT UserID, Username, Password FROM Users WHERE Username=$1", username).Scan(&user.UserID, &user.Username, &user.Password)
	checkInternalServerError(err, w)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	sess := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"username": user.Username, "userid": user.UserID},
		Attrs:  map[string]interface{}{"count": 1},
	})
	session.Add(sess, w)

	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

func (a *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	log.Printf("User %s", s.CAttr("username").(string))
	session.Remove(s, w)
	s = nil

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (a *App) isAuthenticated(w http.ResponseWriter, r *http.Request) {
	authenticated := false

	sess := session.Get(r)

	if sess != nil {
		u := sess.CAttr("username").(string)
		c := sess.Attr("count").(int)

		if c > 0 && len(u) > 0 {
			authenticated = true
		}
	}

	if !authenticated {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func (a *App) setupAuth() {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})
}
