package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/icza/session"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("registerHandler")

	// Render the registration page
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

	//  If the user does not exist, insert the user into the database
	switch {
	case err == sql.ErrNoRows:
		// encrypt the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)

		// insert to database
		_, err = a.db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", username, hashedPassword)
		log.Printf("inserted user to the database")

		checkInternalServerError(err, w)
		// Render the login page after successful registration
		http.ServeFile(w, r, "tmpl/login.html")

	// Redirect to the login page after successful registration
	case err != nil:
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("loginHandler")

	// Render the login page
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

	//password is encrypted
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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

	// Redirect to the list page after successful login
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (a *App) isAuthenticated(w http.ResponseWriter, r *http.Request) {
	log.Printf("isAuthenticated")

	// if the user is authenticated, return true
	authenticated := false

	// Get the current session
	sess := session.Get(r)

	// Check if the session is valid
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
	// Set the session store
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("updateUser")

	// Render the user settings page
	a.isAuthenticated(w, r)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	// Extract user information from the form
	newUsername := r.FormValue("newUsername")
	newPassword := r.FormValue("newPassword")
	confirmPass := r.FormValue("confirmPassword")

	// Get the current session
	s := session.Get(r)
	currentUserID := s.CAttr("userid")

	switch {
	case newUsername != "" && newPassword != "":
		// Check if the new password matches the confirm password
		if newPassword == confirmPass {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			checkInternalServerError(err, w)
			SQL, err := a.db.Prepare(`UPDATE "users" SET username=$1, password=$2 WHERE userid=$3`)
			checkInternalServerError(err, w)
			SQL.Exec(newUsername, hashedPassword, currentUserID)
		} else {
			http.Redirect(w, r, "/", 200)
		}
	case newUsername != "" && newPassword == "":
		// Update the username
		SQL, err := a.db.Prepare(`UPDATE "users" SET username=$1 WHERE userid=$2`)
		checkInternalServerError(err, w)
		SQL.Exec(newUsername, currentUserID)
	case newUsername == "" && newPassword != "":
		// Update the password
		if newPassword == confirmPass {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			checkInternalServerError(err, w)
			SQL, err := a.db.Prepare(`UPDATE "users" SET password=$1 WHERE userid=$2`)
			checkInternalServerError(err, w)
			SQL.Exec(hashedPassword, currentUserID)
		} else {
			http.Redirect(w, r, "/", 200)
		}
	default:
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Redirect to a success page or user settings page
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (a *App) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("deleteUserHandler")

	a.isAuthenticated(w, r)

	sess := session.Get(r)
	sessUserID := sess.CAttr("userid").(int)

	// Extract user information from the form
	deleteUsername := r.FormValue("deleteUsername")
	var deleteUserID int
	err := a.db.QueryRow(`SELECT userid FROM "users" WHERE username=$1`, deleteUsername).Scan(&deleteUserID)
	if err != nil {
		log.Printf("Error retrieving user from the database")
		checkInternalServerError(err, w)
		return
	}

	// Check if the user is deleting their own account
	if deleteUserID == sessUserID {
		a.db.Exec(`DELETE FROM "users" WHERE userid=$1`, deleteUserID)
	} else {
		log.Print("Error matching id")
		http.Redirect(w, r, "/", 200)
	}

	session.Remove(sess, w)
	sess = nil

	// Redirect the user to the login page after successful deletion
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (a *App) createsharinglistHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Create sharing list")
	a.isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	sess := session.Get(r)
	currentuserid := sess.CAttr("userid").(int)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		http.Redirect(w, r, "/list", http.StatusBadRequest)
	}

	listname := r.FormValue("listname")
	log.Println(listname)
	selectedusers := r.Form["user"]

	stmt, err := a.db.Prepare("INSERT INTO custom_sharing_lists (userID, listname) VALUES ($1,$2)")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/list", http.StatusBadRequest)
	}

	_, err = stmt.Exec(currentuserid, listname)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/list", http.StatusBadRequest)
	}

	for _, userIDStr := range selectedusers {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/list", http.StatusBadRequest)
		}
		_, err = a.db.Exec("INSERT INTO user_custom_sharing_lists (userID, listID) VALUES ($1, (SELECT listID FROM custom_sharing_lists WHERE userID = $2 AND listname = $3))", userID, currentuserid, listname)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/list", http.StatusBadRequest)
		}
	}

	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}
