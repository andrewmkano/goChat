package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var connections []*websocket.Conn

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, w.Header())
		connections = append(connections, conn)

		go func(conn *websocket.Conn) {
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
				}
				for _, cox := range connections {
					cox.WriteMessage(msgType, msg)
					println(string(msg))
				}
			}
		}(conn)
	})

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
