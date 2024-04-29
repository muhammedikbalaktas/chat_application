package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/get_user", getUser)
	r.HandleFunc("/", messageHandler)
	r.HandleFunc("/get_contacts", getContacts)
	r.HandleFunc("/list_messages", listMessages)

	go handleMessages()
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
