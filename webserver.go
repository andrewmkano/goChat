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
var connections []*websocket.Conn //Connections stablished

//Creating a struct for the channels
type chatChannel struct {
	channelName string `json:"channelName"` //A name asigned to the channel(Default is always present)
	path        string `json:"path"`        //the path that leads to the channel(Need a handler for this)
}

type chatChannels struct {
	chatChannels []chatChannel `json:"chatChannels"`
}

//Today's objective: Make different channels that people can use to comunicate. i have no clue how to do this
func main() {
	var c chatChannels
	router := mux.NewRouter()
	//Creating Default channel for all the newcomers
	nChannel := createChannel("Default")
	c.chatChannels = append(c.chatChannels, nChannel)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, w.Header())
		connections = append(connections, conn)

		conn.SetCloseHandler(func(code int, text string) error {
			println("%v %v \n", code, text)
			return nil
		})

		go func(conn *websocket.Conn) {
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					println("upgrader error %s\n" + err.Error())
					conn.Close()
					return
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
		IdleTimeout:  5 * time.Minute,
	}

	panic(server.ListenAndServe())
}

func createChannel(channName string) chatChannel {
	chPath := "/" + channName
	return chatChannel{
		channelName: channName,
		path:        chPath,
	}
}
