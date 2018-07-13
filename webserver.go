package main

import (
	"encoding/json"
	"fmt"
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
	changeType  = "CHANGE"
)

var upgrader = websocket.Upgrader{}

//Connections stablished
var connections []*websocket.Conn

var c Channels
var users []User

//Creating a struct for the channels
type Channel struct {
	//name is a name asigned to the channel(Default is always present)
	Name  string `json:"name"`
	Users []User
}

func (cc Channel) String() string {
	return fmt.Sprintf("%s - %v - %v", cc.Name, len(users), users)
}

type Channels []*Channel

func (cl Channels) Debug() {
	for _, ch := range cl {
		fmt.Println(ch)
	}
}

type User struct {
	Connection *websocket.Conn
}

type Message struct {
	Type        string      `json:"Type"`
	Body        interface{} `json:"Body"`
	ChannelName string      `json:"ChannelName"`
}

func init() {
	c = append(c, &Channel{
		Name:  "Default",
		Users: users,
	})

}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/ws", websocketHandler)
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

func addUserToDefault(usr User) {
	log.Println("Adding to default")
	users = append(users, usr)
	c[0].Users = users
}

//Extracts the user connection from the current slice
func getUserConnection(currentChannel string, connection *websocket.Conn) User {
	var usrSwitching User
channelLoop:
	for i := range c {
		if c[i].Name != currentChannel {
			continue
		}
		channelUsers := c[i].Users
		for k := range channelUsers {
			if channelUsers[k].Connection != connection {
				continue
			}
			usrSwitching = channelUsers[k]
			c[i].Users = append(channelUsers[:k], channelUsers[k+1:]...)
			break channelLoop
		}
	}

	return usrSwitching
}

func getChannel(channelName string) *Channel {
	var foundChannel *Channel
	for i := range c {
		if c[i].Name == channelName {
			foundChannel = c[i]
		}
	}
	return foundChannel
}

func UserChannel(usr User) string {
	var userChan string
userLoop:
	for i := range c {

		for j := range c[i].Users {
			if c[i].Users[j].Connection != usr.Connection {
				continue
			}
			userChan = c[i].Name
			break userLoop
		}
	}
	return userChan
}
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	var conn, _ = upgrader.Upgrade(w, r, w.Header())
	connections = append(connections, conn)

	currentUser := User{
		Connection: conn,
	}
	addUserToDefault(currentUser)

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

			messag := Message{}
			err = json.Unmarshal(msg, &messag)

			if err != nil {
				log.Println(err)
				return
			}

			switch messag.Type {
			case "MESSAGE":
				targetChannel := getChannel(messag.ChannelName)
				for _, usr := range targetChannel.Users {
					usr.Connection.WriteJSON(messag)
					log.Print(msg, msgType)
				}
			case "NEW_CHANNEL":

				channelName, ok := messag.Body.(string)
				if !ok {
					println("Not a String dude!")
					continue
				}
				var usersList []User

				c = append(c, &Channel{
					Name:  channelName,
					Users: usersList,
				})
				broadcastChannels()

			case "CHANGE":
				newchannel, ok := messag.Body.(string)
				if !ok {
					println("Not a String dude!")
					continue
				}
				currentChannel := getChannel(UserChannel(routineUser))
				nextChannel := getChannel(newchannel)
				user := getUserConnection(currentChannel.Name, currentUser.Connection)
				nextChannel.Users = append(nextChannel.Users, user)
			}

		}
	}(conn)
}
