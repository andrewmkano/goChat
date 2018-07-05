package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var upgrader = websocket.Upgrader{}
var connections []*websocket.Conn

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "wsocket.html")
	})

	router.HandleFunc("/g", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		connections = append(connections, conn)

		go func(conn *websocket.Conn) {
			for {
				msgType, msg, _ := conn.ReadMessage()

				for _, cox := range connections {
					cox.WriteMessage(msgType, msg)
					log.Println("Msg", cox, string(msg))
				}
			}
		}(conn)
	})

	router.HandleFunc("/l", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			for {
				_, msg, _ := conn.ReadMessage()
				println(string(msg))
			}
		}(conn)
	})

	logCreator := handlers.LoggingHandler(os.Stdout, router)

	server := http.Server{
		Addr:         "127.0.0.1:3000",
		Handler:      logCreator,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	panic(server.ListenAndServe())
}
