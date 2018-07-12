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

const (
	messageType = "MESSAGE"
	channelType = "CHANNEL"
)

var upgrader = websocket.Upgrader{}

//Connections stablished
var connections []*websocket.Conn
var c []Channel
var users []User

//Creating a struct for the channels
type Channel struct {
	//name is a name asigned to the channel(Default is always present)
	Name string `json:"name"`
}

type User struct {
	connection     *websocket.Conn
	currentChannel Channel
}
type Message struct {
	Type string      `json:"Type"`
	Body interface{} `json:"Body"`
}

func init() {
	c = append(c, Channel{
		Name: "Default",
	})

}

func main() {
	router := mux.NewRouter()

	//Websocket connection handler
	router.HandleFunc("/ws", websocketHandler)
	//Serving static files
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

func broadcastChannels() {
	for _, conn := range connections {
		err := notifyChannels(conn)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func notifyChannels(conn *websocket.Conn) error {
	return conn.WriteJSON(Message{
		Type: channelType,
		Body: c,
	})
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	var conn, _ = upgrader.Upgrade(w, r, w.Header())

	connections = append(connections, conn)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := notifyChannels(conn)
	if err != nil {
		log.Fatal(err)
	}

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

			messag := Message{}
			err = json.Unmarshal(msg, &messag)

			if err != nil {
				log.Println(err)
				return
			}

			switch messag.Type {
			case "MESSAGE":
				for _, cox := range connections {
					cox.WriteJSON(messag)
					log.Print(msg, msgType)
				}
			case "NEW_CHANNEL":

				channelName, ok := messag.Body.(string)
				if !ok {
					println("Not a String dude man!")
					continue
				}

				c = append(c, Channel{
					Name: channelName,
				})

				broadcastChannels()
			}

		}
	}(conn)
}
