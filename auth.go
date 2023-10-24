package main

import (
	"database/sql"
	"log"
	"net/http"

	//"unicode"
	//"strings"

	"github.com/icza/session"
	"golang.org/x/crypto/bcrypt"
	//"github.com/asaskevich/govalidator"
)

/*// Validate if the password contains at least one numeric character
func hasNumeric(password string) bool {
    for _, char := range password {
        if unicode.IsDigit(char) {
            return true
        }
    }
    return false
}

// Validate if the password contains at least one special character
func hasSpecialChar(password string) bool {
    specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>/?\\"
    for _, char := range password {
        if strings.ContainsRune(specialChars, char) {
            return true
        }
    }
    return false
}*/

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/register.html")
		return
	}

	// Extract user information from the form
	username := r.FormValue("username")
	password := r.FormValue("password")

	/*// Validate username and password
		passwordValid := govalidator.HasUpperCase(password) &&
	    govalidator.HasLowerCase(password) &&
	    hasNumeric(password) &&
	    hasSpecialChar(password)

		if !passwordValid {
			errorMsg := "Password criteria not met."
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}*/

	// Check if the user already exists in the database
	var user User
	err := a.db.QueryRow(`SELECT username, password FROM "users" WHERE username=$1`, username).Scan(&user.Username, &user.Password) //Scan(&user.Username, &user.Password)

	// ***to do ***      needs mor error checking ] here could using if err or switch case here
	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)
		_, err = a.db.Exec(`INSERT INTO users(username, password, role) VALUES($1, $2, $3)`,
			username, hashedPassword)
		checkInternalServerError(err, w)
	case err != nil:
		http.Error(w, "loi: "+err.Error(), http.StatusBadRequest)
		return
	default:
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
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

	// validate password
	/*
		//simple unencrypted method
		if user.Password != password {
			http.Redirect(w, r, "/login", 301)
			return
		}
	*/

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

// user settings
func (a *App) updatePasswordHandler(w http.ResponseWriter, r *http.Request) {

}
