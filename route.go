package main

import (
	"log"
	"net/http"
)
func (a *App) initalizeRoutes() {
	staticFileDirectory := http.Dir("./statics/")
	staticFileHandler := http.StripPrefix("/statics/", http.FileServer(staticFileDirectory))
	a.Router.PathPrefix("/statics/").Handler(staticFileHandler).Methods("GET")
	a.Router.HandleFunc("/", a.indexHandler).Methods("GET")
	a.Router.HandleFunc("/login", a.loginHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/logout", a.logoutHandler).Methods("GET")
	a.Router.HandleFunc("/register", a.registerHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/list", a.listHandler).Methods("GET")
	a.Router.HandleFunc("/list/{srt:[0-9]+}", a.listHandler).Methods("GET")
	a.Router.HandleFunc("/create", a.createHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/update", a.updateHandler).Methods("POST", "GET")
	a.Router.HandleFunc("/delete", a.deleteHandler).Methods("POST", "GET")

	log.Println("Routes established")
}