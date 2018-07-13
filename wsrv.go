package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func start() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", wsocketHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	logCreator := handlers.LoggingHandler(os.Stdout, router)
	server := http.Server{
		Addr:         "0.0.0.0:3000",
		Handler:      logCreator,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	panic(server.ListenAndServe())
}

func addUserToDefaultChannel(usr user) {
	users = append(users, usr)
	c[0].Users = users
}

func sendChannelsList(conn *websocket.Conn) error {
	var msg struct {
		Type string
		Body interface{}
	}
	msg.Type = channelType
	msg.Body = channels
	return conn.WriteJSON(msg)
}

func broadcastChannelsList() {
	for _, conn := range users {
		err := sendChannelsList(conn.Connection)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func wsocketHandler(w http.ResponseWriter, r *http.Request) {
	var conn, _ = upgrader.Upgrade(w, r, w.Header())
	currentUser := user{
		Connection: conn,
	}
	users = append(users, currentUser)
	addUserToDefaultChannel(currentUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := sendChannelsList(conn)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		println("%v %v \n", code, text)
		return nil
	})

	go func(conn *websocket.Conn) {
		routineUser := User{
			Connection: conn,
		}

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				println("upgrader error %s\n" + err.Error())
				conn.Close()
				return
			}

			var messg struct {
				Type string
				Body interface{}
			}

			err = json.Unmarshal(msg, &messg)

			if err != nil {
				log.Println(err)
				return
			}

			switch messg.Type {
			case "MESSAGE":
				messageText, ok := messg.Body.(string)
				if !ok {
					println("Not a String dude!")
					continue
				}
				receiveAndSendMessage(messageText, routineUser.Connection)
			case "NEW_CHANNEL":
				channelName, ok := messg.Body.(string)
				if !ok {
					println("Not a String dude!")
					continue
				}
				createAndAddChannel(channelName, routineUser.Connection)
			case "CHANGE":
				newchannel, ok := messg.Body.(string)
				if !ok {
					println("Not a String dude!")
					continue
				}
				UpdateUserChannel(newchannel, routineUser.Connection)
			}

		}
	}(conn)
}
