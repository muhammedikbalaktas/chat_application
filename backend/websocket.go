package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{

		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	userConnections = make(map[int]*websocket.Conn)
	broadcast       = make(chan []byte)

	mutex sync.Mutex
)

type Message struct {
	Text        string `json:"text" validate:"required"`
	Reciever_id int    `json:"reciever_id" validate:"required"`
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get the value of the "token" query parameter
	token := query.Get("token")
	userId, err := parseToken(token)
	if err != nil {
		fmt.Println(token)
		fmt.Println(err)
		fmt.Println("error on parsing token")
		return
	}
	fmt.Println("someone connected with userId ", userId)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	userConnections[userId] = conn

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Read error:", err)
			break
		}

		broadcast <- message
	}

	removeConnection(userId)
	fmt.Println("user exited")
}
func removeConnection(userId int) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(userConnections, userId)

}

func handleMessages() {

	for {

		var mess Message
		message := <-broadcast
		err := json.Unmarshal(message, &mess)
		fmt.Println("recieved message is", string(message))
		if err != nil {
			fmt.Println("invalid message")
			fmt.Println(err)
			continue
		}
		validator := validator.New()
		err = validator.Struct(mess)
		if err != nil {
			fmt.Println("invalid json from validator")
			continue
		}
		for id, conn := range userConnections {

			if mess.Reciever_id == id {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					fmt.Println("write error ", err)
					conn.Close()
					removeConnection(id)
				}
				fmt.Println(string(message))
				addMessageToDatabase(id, mess)
			}
			addMessageToDatabase(id, mess)

		}
	}
}
func parseToken(tokenString string) (int, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
func addMessageToDatabase(senderId int, message Message) {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error opening databases")
		removeConnection(senderId)
		return

	}
	query := "insert into messages (sender,reciever,time,message) values(?,?, now(),?)"
	result, err := db.Exec(query, senderId, message.Reciever_id, message.Text)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on inserting data")
		removeConnection(senderId)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		fmt.Println("error on adding data")
		removeConnection(senderId)
		return
	}
	fmt.Println("message added succesfully")
}
