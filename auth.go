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
		http.ServeFile(w, r, "templates/register.html")
		return
	}

	//user information
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	// Check existence of user
	var user User
	err := a.db.QueryRow(`SELECT username, password, role FROM "users" WHERE username=$1`, username).Scan(&user.Username, &user.Password, &user.Role)
	switch {
	// user is availabl
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)
		// insert to database
		_, err = a.db.Exec(`INSERT INTO "users"(username, password, role) VALUES($1, $2, $3)`,
			username, hashedPassword, role)
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
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	//user info from the submitted form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// query database to get match username
	var user User
	err := a.db.QueryRow(`SELECT userID, username, password FROM "users" WHERE username=$1`, username).Scan(&user.UserID, &user.Username, &user.Password)
	checkInternalServerError(err, w)

	//password is encrypted
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

	// get the current session variables
	s := session.Get(r)
	log.Printf("User %s", s.CAttr("username").(string))
	session.Remove(s, w)
	s = nil

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (a *App) isAuthenticated(w http.ResponseWriter, r *http.Request) {
	log.Printf("Authrntication")
	authenticated := false

	sess := session.Get(r)

	if sess != nil {
		u := sess.CAttr("username").(string)
		c := sess.Attr("count").(int)

		//authentication check for the current user
		if c > 0 && len(u) > 0 {
			authenticated = true
		}
	}

	if !authenticated {
		log.Printf("Authentication failed")
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func (a *App) setupAuth() {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})
}
