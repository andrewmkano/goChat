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
	//ChannelName is a name asigned to the channel(Default is always present)
	ChannelName string `json:"ChannelName"`
	//Path is the Path that leads to the channel(Need a handler for this)
	Path string `json:"Path"`
}

// type chatChannels []chatChannel

type chatChannels struct {
	chatChannels []chatChannel `json:"chatChannels"`
}

//Today's objective: Make different channels that people can use to comunicate. i have no clue how to do this
func main() {
	var c chatChannels
	router := mux.NewRouter()

	//Creating Default channel for all the newcomers
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, w.Header())
		connections = append(connections, conn)
		nChannel := createChannel("Default", w, r)
		c.chatChannels = append(c.chatChannels, nChannel)
		// jsonFile, err := json.Marshal(nChannel)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		conn.WriteJSON(nChannel)
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

func createChannel(channName string, w http.ResponseWriter, r *http.Request) chatChannel {
	chPath := "/" + channName
	return chatChannel{
		ChannelName: channName,
		Path:        chPath,
	}
}
