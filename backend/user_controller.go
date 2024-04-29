package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbInfo    = "<your_username>:<your_password>@tcp(localhost:3306)/<database_name>"
	secretKey = []byte("mySecretKey")
)

func getUser(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is OPTIONS (preflight)
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

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errorMessage := map[string]string{"error": err.Error()}
		jsonResponse, _ := json.Marshal(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		fmt.Println(err)
		return
	}

	validator := validator.New()
	if err := validator.Struct(user); err != nil {
		http.Error(w, "Failed to validate JSON", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		fmt.Println("error opening database")
		fmt.Println(err)
		return
	}
	defer db.Close() // Close the database connection when done

	query := "SELECT id FROM users WHERE username=? AND password=?"

	result, err := db.Query(query, user.Username, user.Password)
	if err != nil {
		fmt.Println("error finding user")
		fmt.Println(err)
		return
	}
	defer result.Close() // Close the result set when done

	var (
		id       int
		rowCount int
	)

	for result.Next() {
		rowCount++
		err = result.Scan(&id)
		if err != nil {
			fmt.Println("error in getting data")
		}
	}
	if rowCount == 0 {
		fmt.Println("no such user")
		return
	}

	token, err := generateToken(id)
	if err != nil {
		fmt.Println("error generating token")
		fmt.Println(err)
		return
	}
	var mapped = map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(mapped)
	if err != nil {
		fmt.Println("error marshalling token")
		return
	}
	w.Write(response)
}

func generateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})
	tokeString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokeString, nil
}
