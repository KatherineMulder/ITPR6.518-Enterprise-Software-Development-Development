package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/icza/session"
	"golang.org/x/crypto/bcrypt"
)

//// The registerHandler handles user registration.
func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {

 	// Serve the registration form if the request method is not POST
	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/register.html")
		return
	}

	// Extract user information from the form
	username := r.FormValue("username")
	password := r.FormValue("password")


	// Check if the user already exists in the database
	var user User
	err := a.db.QueryRow(`SELECT username, password FROM "users" WHERE username=$1`, username).Scan(&user.Username, &user.Password)
	switch {
	// user is available
	case err == sql.ErrNoRows:

		// User is not found, so we hash the password and insert it into the database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)

		// insert to database
		//a.db.Exec(`SELECT setval(pg_get_serial_sequence('users', 'userid'), coalesce(max(id),0) + 1, false) FROM t1;`)
		_, err = a.db.Exec(`INSERT INTO "users"(username, password) VALUES($1, $2)`,
			username, hashedPassword)
		checkInternalServerError(err, w)
	case err != nil:
		 // An error occurred during the database query
		http.Error(w, "loi: "+err.Error(), http.StatusBadRequest)
		return
	default:
		 // User already exists, redirect to the login page
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
	// Redirect to the login page after registration
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Method %s", r.Method)
	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/login.html")
		return
	}

	//user info from the submitted form
	username := r.FormValue("usrname")
	log.Println(username)
	password := r.FormValue("psw")
	log.Println(password)

	// query database to get match username
	var user User
	err := a.db.QueryRow(`SELECT userID, username, password FROM "users" WHERE username=$1`, username).Scan(&user.UserID, &user.Username, &user.Password)
	checkInternalServerError(err, w)

	//password is encrypted
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	log.Println(user.Password)
	if err != nil {
		if password == user.Password {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			checkInternalServerError(err, w)
			// insert to database
			_, err = a.db.Exec(`UPDATE "users" SET password=$1 WHERE username=$2`, hashedPassword, username)
			checkInternalServerError(err, w)
		} else {
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}
	}

 	// Create and store a session
	sess := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"username": user.Username, "userid": user.UserID},
		Attrs:  map[string]interface{}{"count": 1},
	})
	session.Add(sess, w)

	// Redirect to the list page after successful login
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

func (a *App) logoutHandler(w http.ResponseWriter, r *http.Request) {

	// get the current session variables
	s := session.Get(r)
	log.Printf("User %s", s.CAttr("username").(string))
	session.Remove(s, w)
	s = nil

	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// The isAuthenticated function checks if the user is authenticated.
func (a *App) isAuthenticated(w http.ResponseWriter, r *http.Request) {
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
 	// Redirect to the login page if the user is not authenticated
	if !authenticated {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func (a *App) setupAuth() {
	
	// Set up session management
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})
}
