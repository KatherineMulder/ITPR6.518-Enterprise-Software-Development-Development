package main

import (
	"database/sql"
	"log"
	"net/http"
	//"html/template" for cookies

	"github.com/icza/session"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/register.html")
		return
	}

	// Extract user information from the form
	username := r.FormValue("username")
	password := r.FormValue("password")

	
	// Check if the user already exists in the database
	var user User
	err := a.db.QueryRow(`SELECT username, password FROM "users" WHERE username=$1`, username).Scan(&user.Username)  //Scan(&user.Username, &user.Password)


	// ***to do ***      needs mor error checking ] here could using if err or switch case here  
	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)

		// insert to database
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

/*//password hashing
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(hash), err
}

func CheckPasswordWithHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}*/

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method %s", r.Method)


	if r.Method != "POST" {
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

	//user info from the submitted form
	username := r.FormValue("username")
	log.Println(username)

	password := r.FormValue("psw")
	log.Println(password)


	// query database to get match username
	var user User
	err = a.db.QueryRow("SELECT username, password FROM users WHERE username=$1", username).Scan(&user.Username, &user.Password)

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
		CAttrs: map[string]interface{}{"username": user.Username},
		Attrs:  map[string]interface{}{"count": 1},
	})
	session.Add(sess, w)

	// Redirect to the list page after successful login
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)

        return
    }
}


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

func (a *App) logoutHandler(w http.ResponseWriter, r *http.Request) {

	// get the current session variables
	s := session.Get(r)
	log.Printf("User %s", s.CAttr("username").(string))
	session.Remove(s, w)
	s = nil

	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// authentication handlers using the sessions
func (a *App) setupAuth() {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})
}
