package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/icza/session"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("registerHandler")

	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/register.html")
		return
	}

	// Extract user information from the form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if the user already exists in the database
	var user User
	err := a.db.QueryRow("SELECT username, password FROM users WHERE username=$1", username).Scan(&user.Username, &user.Password)
	log.Printf("User %s", user.Username)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)

		// insert to database
		_, err = a.db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", username, hashedPassword)
		log.Printf("inserted user to the database")

		checkInternalServerError(err, w)
		// Render the login page after successful registration
		http.ServeFile(w, r, "tmpl/login.html")

	case err != nil:
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("loginHandler")

	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/login.html")
		return
	}

	// grab user info from the submitted form
	username := r.FormValue("usrname")
	password := r.FormValue("psw")

	// query database to get match username
	var user User
	err := a.db.QueryRow("SELECT userid, username, password FROM users WHERE username=$1",
		username).Scan(&user.UserID, &user.Username, &user.Password)
	checkInternalServerError(err, w)
	log.Println(user)

	//password is encrypted
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	log.Println(err)
	if err != nil {
		if user.Password == password {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			checkInternalServerError(err, w)
			_, err = a.db.Exec(`UPDATE users SET password=$1 WHERE username=$2`,
				hashedPassword, username)
			checkInternalServerError(err, w)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	// Successful login. New session with initial constant and variable attributes
	sess := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"username": user.Username, "userid": user.UserID},
		Attrs:  map[string]interface{}{"count": 1},
	})
	session.Add(sess, w)

	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (a *App) isAuthenticated(w http.ResponseWriter, r *http.Request) {
	log.Printf("isAuthenticated")

	authenticated := false

	// Get the current session
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
	log.Printf("logoutHandler")

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

func (a *App) updateUserSetting(w http.ResponseWriter, r *http.Request) {
	log.Printf("updateUserSetting")

	if r.Method != "POST" {
		// Serve a form to allow users to update their settings
		http.ServeFile(w, r, "tmpl/update_settings.html")
		return
	}

	// Extract user information from the form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if a new username is provided
	if username != "" {
		// Update the user's username in the database
		_, err := a.db.Exec("UPDATE users SET username=$1 WHERE username=$2", username, session.Get(r).CAttr("username").(string))
		if err != nil {
			http.Error(w, "Error updating username: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	log.Printf("Updated username")

	// Check if a new password is provided
	if password != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			checkInternalServerError(err, w)
			return
		}

		// Update the user's password in the database
		_, err = a.db.Exec("UPDATE users SET password=$1 WHERE username=$2", hashedPassword, username)
		if err != nil {
			http.Error(w, "Error updating password: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	log.Printf("Updated password")

	// Redirect to a success page or user settings page
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (a *App) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("deleteUserHandler")

	a.isAuthenticated(w, r)

	sess := session.Get(r)
	userID := sess.CAttr("userid").(int)

	// Check if the user exists
	var deletedUser User // Replace 'User' with your user struct type
	err := a.db.QueryRow("SELECT * FROM users WHERE userid = $1", userID).Scan(&deletedUser.UserID, &deletedUser.Username, &deletedUser.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_, err = a.db.Exec("DELETE FROM users WHERE userid = $1", userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the login page after successful deletion
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
