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

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("updateUser")
	a.isAuthenticated(w, r)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	// Extract user information from the form
	newUsername := r.FormValue("newUsername")
	newPassword := r.FormValue("newPassword")
	confirmPass := r.FormValue("confirmPassword")

	s := session.Get(r)
	currentUserID := s.CAttr("userid")

	switch {
	case newUsername != "" && newPassword != "":
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
		SQL, err := a.db.Prepare(`UPDATE "users" SET username=$1 WHERE userid=$2`)
		checkInternalServerError(err, w)
		SQL.Exec(newUsername, currentUserID)
	case newUsername == "" && newPassword != "":
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

	deleteUsername := r.FormValue("deleteUsername")
	log.Println(deleteUsername)
	var deleteUserID int
	err := a.db.QueryRow(`SELECT userid FROM "users" WHERE username=$1`, deleteUsername).Scan(&deleteUserID)
	log.Println(err)
	if err != nil {
		log.Printf("Error retrieving user from the database")
		checkInternalServerError(err, w)
		return
	}

	if deleteUserID == sessUserID {
		a.db.Exec(`DELETE FROM "users" WHERE userid=$1`, deleteUserID)
	} else {
		log.Print("Error matching id")
		http.Redirect(w, r, "/", 200)
	}

	session.Remove(sess, w)
	sess = nil

	// Redirect the user to the login page after successful deletion
<<<<<<< Updated upstream
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
=======
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
>>>>>>> Stashed changes
