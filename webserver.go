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
	}
	panic(server.ListenAndServe())
}

//Updates the list of channels for everyone
func broadcastChannels() {
	for _, conn := range connections {
		err := notifyChannels(conn)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//Sends a message that contains the list of channels
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
	//Iterate over the channels

channelLoop:
	for i := range c {
		if c[i].Name != currentChannel {
			continue
		}
		//Look for the current channel of the user
		channelUsers := c[i].Users
		//Iterate over the users on that channel
		for k := range channelUsers {
			//Look for the user connection
			if channelUsers[k].Connection != connection {
				continue
			}

			usrSwitching = channelUsers[k]
			c[i].Users = append(channelUsers[:k], channelUsers[k+1:]...)
			break channelLoop
		}
	}

	// c.Debug()
	return usrSwitching
}

func getNextChannel(newChannel string, connection *websocket.Conn) []User {
	var channelUsers []User
	channelUsers = nil

	for i := range c {
		//Look for the next channel of the user
		if c[i].Name == newChannel {
			channelUsers = c[i].Users
			return channelUsers
		}
	}
	return channelUsers
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

//Switches the user from a channel to another
func SwitchChannel(user User, newChannel string) {
	//Finds the channel
	nextChannelUsers := getNextChannel(newChannel, user.Connection)
	//Moves the user into the channel
	nextChannelUsers = append(nextChannelUsers, user)
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
				println("Im on NEW CHANNEL")
				println(c[0].Name, " Users: ", len(c[0].Users))
				println(c[1].Name, " Users: ", len(c[1].Users))
				println("#########")
			case "CHANGE":
				newchannel, ok := messag.Body.(string)
				if !ok {
					println("Not a String dude!")
					continue
				}
				println("Im on CHANGE")
				println("ChannelUsers[Default]: ", len(c[0].Users))
				println("#########")
				currentChannel := getChannel(UserChannel(routineUser))
				nextChannel := getChannel(newchannel)
				user := getUserConnection(currentChannel.Name, currentUser.Connection)
				println("Before")
				println(currentChannel.Name, "Users on this channel: ", len(currentChannel.Users), user.Connection)
				println(nextChannel.Name, "Users on this channel: ", len(nextChannel.Users))
				println("Captured User", user.Connection)
				println("#########")
				nextChannel.Users = append(nextChannel.Users, user)
				// SwitchChannel(user, newchannel)
				println("After")
				println(currentChannel.Name, "Users on this channel: ", len(currentChannel.Users))
				println(nextChannel.Name, "Users on this channel: ", len(nextChannel.Users), nextChannel.Users)
				println("#########")
			}

		}
	}(conn)
}
