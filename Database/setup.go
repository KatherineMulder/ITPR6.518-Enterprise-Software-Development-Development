package setup

import (
	"database/sql"
	"fmt"
	"log" 

	_ "github.com/lib/pq" 

)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "0518"
	dbname   = "enterpriseNotes"

	// HTTP 

)

//create a database connection

func createDataBase() error {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return err
    }

    defer db.Close()

    err = db.Ping()
    if err != nil {
        return err
    }

    fmt.Println("Connected successfully")

    // Create database query
    sqlStatement := `SELECT nameid, firstname, lastname, gender, city, job
	 FROM tblemployees WHERE nameid >= $1 AND nameid <= $2;`

    rows, err := db.Query(sqlStatement, 3, 5)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var nameid int
        var firstname string
        var lastname string
        var gender string
        var city string
        var job string

        switch err = rows.Scan(&nameid, &firstname, &lastname, &gender, &city, &job); err {
        case sql.ErrNoRows:
            fmt.Println("No rows were returned!")
        case nil:
            fmt.Println(nameid, firstname, lastname, job)
        default:
            return err
        }
    }

    // Get any error encountered during iteration
    err = rows.Err()
    if err != nil {
        return err
    }

    // Return nil if everything is successful
    return nil
}
