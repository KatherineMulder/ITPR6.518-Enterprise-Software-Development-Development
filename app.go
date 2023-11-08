package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Constants for database connection
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "EnterpriseNotes"
)

// Variables for database connection
var (
	err  error
	wait time.Duration
)

// App represents the main application.
type App struct {
	Router   *mux.Router
	db       *sql.DB
	bindport string
	//username string
}

// Initialize initializes the application with predefined configuration
func (a *App) Initialize() {

	a.bindport = "8080"

	// check if the port is set in the environment variables
	tempport := os.Getenv("PORT")
	if tempport != "" {
		a.bindport = tempport
	}

	// check if the port is set in the command line arguments
	if len(os.Args) > 1 {
		s := os.Args[1]

		if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			log.Printf("Connected to port %s", s)
			a.bindport = s
		}
	}

	//database connection
	psqInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	log.Printf("Connecting to PostgresSQL Server")
	log.Println(psqInfo)
	db, err := sql.Open("pgx", psqInfo)
	a.db = db

	//a.importData()
	if err != nil {
		log.Println("Either missing github.com/lib/pq or Invalid DB arguements")
		log.Fatal(err)
	}

	// test connection
	err = a.db.Ping()
	if err != nil {
		log.Fatal("Connection to DB failed: ", err)
	}

	log.Println("Connection to DB successful")

	//check if the imported folder exists
	_, err = os.Stat("./imported")
	if os.IsNotExist(err) {
		log.Println("Importing data")
		a.importData()
	}

	//set some defaults for the authentication to also support HTTP and HTTPS
	a.setupAuth()

	//set up the routes
	a.Router = mux.NewRouter()
	a.initalizeRoutes()

}

// initalizeRoutes initializes the application routes.
func (a *App) initalizeRoutes() {

	//static file routes
	staticFileDirectory := http.Dir("./statics/")
	staticFileHandler := http.StripPrefix("/statics/", http.FileServer(staticFileDirectory))
	a.Router.PathPrefix("/statics/").Handler(staticFileHandler).Methods("GET")

	a.Router.HandleFunc("/", a.indexHandler).Methods("GET")

	//user settings routes
	a.Router.HandleFunc("/login", a.loginHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/logout", a.logoutHandler).Methods("GET")
	a.Router.HandleFunc("/register", a.registerHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/updateUser", a.updateUser).Methods("POST", "GET")
	a.Router.HandleFunc("/deleteUser", a.deleteUserHandler).Methods("POST")
	a.Router.HandleFunc("/createsharinglist", a.createsharinglistHandler).Methods("POST", "GET")

	//note routes
	a.Router.HandleFunc("/list", a.listHandler).Methods("GET")
	a.Router.HandleFunc("/list/{srt:[0-9]+}", a.listHandler).Methods("GET")
	a.Router.HandleFunc("/create", a.createHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/update", a.updateHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/delete", a.deleteHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/sharenote", a.shareNoteHandler).Methods("POST", "GET")

	//search routes
	a.Router.HandleFunc("/searchnotes", a.searchNotesHandler).Methods("GET")
	a.Router.HandleFunc("/getdelegations", a.getdelegationsHandler).Methods("GET")
	a.Router.HandleFunc("/getsharelist", a.getShareListHandler).Methods("GET")
	a.Router.HandleFunc("/getcustomsharinglists", a.getCustomSharingListsHandler).Methods("GET")

	log.Println("Routes established")
}

// Run starts the application.
func (a *App) Run(addr string) {

	// check if the port is set in the environment variables
	if addr != "" {
		a.bindport = addr
	}

	// check if the port is set in the command line arguments
	ip := GetOutboundIP()
	log.Println(ip)
	log.Println(a.bindport)
	log.Printf("Starting EnterpriseNotes via HTTP Services at http://%s:%s", ip, a.bindport)

	// set timeouts so that a slow or malicious client doesn't
	srv := &http.Server{
		Addr: ip + ":" + a.bindport,

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}

	// run our server in a goroutine so that it doesn't block.
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// create a channel to listen for OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	log.Println("Shutting down Web Service")
	srv.Shutdown(ctx)
	log.Println("Disconnecting from DB")
	a.db.Close()
	log.Println("Exiting Program")
	os.Exit(0)
}
