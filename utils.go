package main

import (
	"log"
	"net"
	"net/http"
)

// check if there is an internal server error
func checkInternalServerError(err error, w http.ResponseWriter) {
	log.Printf("Checking for internal server errors")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// check the getoutboundip address of the server
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}


