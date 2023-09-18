package database

import (
	"database/sql"
	"fmt"
	"log" 

	_ "github.com/lib/pq"  // Interface to PostgreSQL library

    //"enterpriseNotes/assignment"

)
// PostgreSQl configuration
const (
	host     = "localhost" 
	port     = 5432
	user     = "postgres"
	password = "0518"
	dbname   = "enterpriseNotes"
)

//create a database connection
var conn *pgx.Conn
//urlExample := "postgres://postgres:postgres@localhost:5432/todo"

func createDataBase() error {
    //do we need this?
    //psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Invalid DB arguments, or github.com/lib/pq not installed")
	}

	defer db.Close() //Ensure connection is always closed once done

	// Ping database 
	err = db.Ping()
	if err != nil {
		log.Fatal("Connection to specified database failed: ", err)

	}

	fmt.Println("Connected successfully")


	// Create database query here. Note $1/$2 variables which are initialized later in db.query as 3 & 5 respectively
	// Other variables can be added as required ($4, $5, $6...)
	sqlStatement := //`SELECT userid, firstname, lastname, email
	 //FROM tblemployees WHERE nameid >= $1 AND nameid <= $2;`

	rows, err := db.Query(sqlStatement, 3, 5) // $1 and $2 set here. Note sqlStatement could be replaced with literal string
	if err != nil {
		log.Fatal(err)
		fmt.Println("An error occurred when querying data!")
	}
	defer rows.Close() // Make sure we free memory resources when done

	for rows.Next() {

		var userid int
		var firstname string
		var lastname string
		var email string
		//password?

		switch err = rows.Scan(&userid, &firstname, &lastname, &email); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			fmt.Println(userid, firstname, lastname, email)
		default:
			fmt.Println("SQL query error occurred: ")
			panic(err)
		}
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)

	}
}

