package main

import (
	"net/http"
	"os"
	"time"

	"html/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var connections []*websocket.Conn

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("public/index.html")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		// http.ServeFile(w, r, "public/index.html")
	})
	router.Handle("/..", http.FileServer(http.Dir("public/styles.css")))

	router.HandleFunc("/.", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/app.js")
	})

	router.HandleFunc("/g", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
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

	logCreator := handlers.LoggingHandler(os.Stdout, router)

	server := http.Server{
		Addr:         "0.0.0.0:3000",
		Handler:      logCreator,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	panic(server.ListenAndServe())
}
