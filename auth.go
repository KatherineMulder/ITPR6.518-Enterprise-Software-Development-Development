package main

import (
	"database/sql"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"html/template"
	
	"strconv"
    "github.com/gorilla/mux"

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
    err := a.db.QueryRow("SELECT username, password FROM users WHERE username=$1", username).Scan(&user.Username, &user.Password)

    switch {
    case err == sql.ErrNoRows:
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        checkInternalServerError(err, w)
        _, err = a.db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", username, hashedPassword)
        checkInternalServerError(err, w)
		
        // Registration successful message
        data := struct {
            Message string
        }{
            Message: "Registration successful. You can now log in.",
        }

        tmpl, err := template.ParseFiles("tmpl/register.html")
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        tmpl.Execute(w, data)  // Execute the template with the message

    case err != nil:
        http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
        return
    default:
        http.Redirect(w, r, "/login", http.StatusMovedPermanently)
    }
}



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

// === CRUD functions for the user ==== //
// reference: https://www.honeybadger.io/blog/how-to-create-crud-application-with-golang-and-mysql/

func (a *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("createUserHandler")
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse JSON data from the request body
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }

    // Check if the user already exists in the database
    var existingUser User
    err := a.db.QueryRow("SELECT username FROM users WHERE username = $1", user.Username).Scan(&existingUser.Username)

    if err == sql.ErrNoRows {
        // User doesn't exist, proceed with registration
        hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

        if hashErr != nil {
            http.Error(w, "Failed to hash the password", http.StatusInternalServerError)
            return
        }

        // Insert the user into the database
        query := "INSERT INTO users (username, password) VALUES ($1, $2)"
        _, insertErr := a.db.Exec(query, user.Username, string(hashedPassword))
        if insertErr != nil {
            http.Error(w, "Failed to insert user into the database", http.StatusInternalServerError)
            return
        }

        // Registration successful
        w.WriteHeader(http.StatusCreated)
		log.Printf("User created successfully")
        fmt.Fprintln(w, "User created successfully")
    } else if err != nil {
        http.Error(w, "Error checking user existence: "+err.Error(), http.StatusInternalServerError)
        return
    } else {
        // User with the same username already exists
        http.Error(w, "Username already exists", http.StatusConflict)
    }
}


func CreateUser(db *sql.DB, username, password string) error {
    query := "INSERT INTO users (username, password) VALUES ($1, $2)"
    _, err := db.Exec(query, username, password)
    return err
}


func (a *App) getUserHandler(w http.ResponseWriter, r *http.Request) {
    // Check the HTTP method
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the 'id' parameter from the URL
    vars := mux.Vars(r)
    idStr := vars["id"]

    // Convert 'id' to an integer
    userID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
        return
    }

    // Call the GetUser function to fetch the user data from the database
    user, err := GetUser(a.db, userID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Convert the user object to JSON and send it in the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func GetUser(db *sql.DB, userID int) (User, error) {
    var user User
    query := "SELECT userID, username, password FROM users WHERE userID = $1"

    err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.Password)
    if err != nil {
        return User{}, err
    }

    return user, nil
}

func (a *App) updateUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the 'id' parameter from the URL
    vars := mux.Vars(r)
    idStr := vars["id"]

    // Convert 'id' to an integer
    userID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
        return
    }

    // Parse JSON data from the request body
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }

    // Call the UpdateUser function to update the user data in the database
    err = UpdateUser(a.db, userID, user.Username, user.Password)
    if err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "User updated successfully")
}

func UpdateUser(db *sql.DB, userID int, username string, password string) error {
    // Write the SQL query to update the user's data in the database
    query := "UPDATE users SET username = $2, password = $3 WHERE userID = $1"
    
    _, err := db.Exec(query, userID, username, password)
    if err != nil {
        return err
    }
    
    return nil
}


func (a *App) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

    // Return a JSON response with the deleted user
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(deletedUser)
}