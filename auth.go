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


func (a *App) updateUsernameHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    newUsername := r.FormValue("newUsername")
    a.updateUserSettings(w, r, newUsername, "") // Empty newPassword as we're not updating the password
}

func (a *App) updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    newPassword := r.FormValue("newPassword")
    confirmPassword := r.FormValue("confirmPassword")

    if newPassword != confirmPassword {
        http.Error(w, "Passwords do not match", http.StatusBadRequest)
        return
    }

    a.updateUserSettings(w, r, "", newPassword) // Empty newUsername as we're not updating the username
}


func (a *App) updateUserSettings(w http.ResponseWriter, r *http.Request, newUsername, newPassword string) {
    a.isAuthenticated(w, r)

    sess := session.Get(r)
    userID := sess.CAttr("userid").(int)

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash the new password", http.StatusInternalServerError)
        return
    }

    _, err = a.db.Exec("UPDATE users SET username = $1, password = $2 WHERE userid = $3", newUsername, hashedPassword, userID)
    if err != nil {
        http.Error(w, "Failed to update user settings", http.StatusInternalServerError)
        return
    }

    updatedSession := session.NewSessionOptions(&session.SessOptions{
        CAttrs: map[string]interface{}{"username": newUsername, "userid": userID},
    })

    session.Remove(sess, w)
    session.Add(updatedSession, w)

    http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (a *App) updateUserSettingsHandler(w http.ResponseWriter, r *http.Request) {
    a.isAuthenticated(w, r)

    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    newUsername := r.FormValue("newUsername")
    newPassword := r.FormValue("newPassword")
    confirmPassword := r.FormValue("confirmPassword")

    if newPassword != confirmPassword {
        http.Error(w, "Passwords do not match", http.StatusBadRequest)
        return
    }

    a.updateUserSettings(w, r, newUsername, newPassword)
}

func (a *App) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	a.isAuthenticated(w, r)

	sess := session.Get(r)
	userID := sess.CAttr("userid").(int)

	_, err := a.db.Exec("DELETE FROM users WHERE userid = $1", userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
   
}




