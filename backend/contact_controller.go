package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Token struct {
	Data string `json:"token" validator:"required"`
}
type Contacts struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// Respond to preflight requests with appropriate CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set CORS headers for the actual request
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var token Token
	contacts := make([]Contacts, 0)

	// Decode request body into token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		errorMessage := map[string]string{"error": err.Error()}
		jsonResponse, _ := json.Marshal(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		fmt.Println(err)
		return
	}

	// Validate token
	validator := validator.New()
	if err := validator.Struct(token); err != nil {
		http.Error(w, "Failed to validate JSON", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Open database connection
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		fmt.Println("error opening database")
		fmt.Println(err)
		return
	}
	defer db.Close() // Close the database connection when done

	// Parse user ID from token
	userId, err := parseToken(token.Data)
	if err != nil {
		fmt.Println("error on parsing token")
		fmt.Println(err)
		return
	}

	// Query contacts for the user
	// 	SELECT id, username
	// FROM users
	// WHERE id IN (SELECT contact_id
	//              FROM users
	//              INNER JOIN contacts ON users.id = contacts.owner_id
	//              WHERE users.id = 1);
	query := "SELECT id, username FROM users WHERE id IN  " +
		"(SELECT contact_id FROM users INNER JOIN contacts ON users.id = contacts.owner_id WHERE users.id = ?)"
	rows, err := db.Query(query, userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close() // Close the result set when done

	// Iterate over query results and populate contacts slice
	for rows.Next() {
		var cnt Contacts
		err := rows.Scan(&cnt.Id, &cnt.Username)
		if err != nil {
			fmt.Println("error on getting rows in database")
			fmt.Println(err)
			return
		}
		contacts = append(contacts, cnt)
	}

	// Respond with contacts in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(contacts)
	if err != nil {
		fmt.Println("error marshalling contacts")
		return
	}
	w.Write(response)
}
