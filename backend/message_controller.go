package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserInfo struct {
	SenderToken string `json:"token" validator:"required"`
	RecieverId  int    `json:"reciever_id" validator:"required"`
}
type Messages struct {
	SenderId   int    `json:"sender_id" validator:"required"`
	RecieverId int    `json:"reciever_id" validator:"required"`
	Time       string `json:"time" validator:"required"`
	Text       string `json:"text" validator:"required"`
}

func listMessages(w http.ResponseWriter, r *http.Request) {
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
	var userInfo UserInfo
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		errorMessage := map[string]string{"error": err.Error()}
		jsonResponse, _ := json.Marshal(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		fmt.Println(err)
		return
	}
	query := "select sender, reciever, time, message from messages " +
		"where (sender=? and reciever=?) or (sender=? and reciever=?) " +
		"order by time asc"
	sql, err := sql.Open("mysql", dbInfo)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error in databases")
		return
	}
	senderId, err := parseToken(userInfo.SenderToken)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on parsing token")
		return
	}

	rows, err := sql.Query(query, senderId, userInfo.RecieverId, userInfo.RecieverId, senderId)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error getting rows")
		return
	}
	var messages = make([]Messages, 0)
	for rows.Next() {
		var mess Messages
		err = rows.Scan(&mess.SenderId, &mess.RecieverId, &mess.Time, &mess.Text)
		if err != nil {
			fmt.Println(err)
			fmt.Println("error on getting messages")
			return
		}
		messages = append(messages, mess)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("error marshalling contacts")
		return
	}
	w.Write(response)
}
