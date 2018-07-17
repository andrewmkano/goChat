package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type chatServer struct {
	users    []user
	channels []channel
}

type message struct {
	Username string
	Text     string
}

type user struct {
	Username   string
	Connection *websocket.Conn
}

type channel struct {
	ChannelName string
	Users       []user
	Messages    []message
}

func (cs *chatServer) receiveAndSendMessage(msg string, userConn *websocket.Conn) {
	receivedMessage := message{
		Text: msg,
	}
searchLoop:
	for _, ch := range cs.channels {
		for j := range ch.Users {
			if userConn != ch.Users[j].Connection {
				continue
			}
			ch.Messages = append(ch.Messages, receivedMessage)
			cs.emitMessages(ch, receivedMessage)
			break searchLoop
		}

	}
}

func (cs *chatServer) addUserToDefaultChannel(usr user) {
	cs.channels[0].Users = append(cs.channels[0].Users, usr)
}

func (cs *chatServer) findChannel(channelName string) *channel {
	var capturedChannel *channel
	for _, ch := range cs.channels {
		if ch.ChannelName != channelName {
			continue
		}
		capturedChannel = &ch
		return capturedChannel
	}
	return capturedChannel
}

func (cs *chatServer) findUserChannel(conn *websocket.Conn) *channel {
	var capturedUserChannel *channel
	for i := range cs.channels {
		for j := range cs.channels[i].Users {
			if cs.channels[i].Users[j].Connection != conn {
				continue
			}
			capturedUserChannel = &cs.channels[i]
			return capturedUserChannel
		}
	}
	return capturedUserChannel
}

func (cs *chatServer) emitMessages(ch channel, msg message) {
	for _, usr := range ch.Users {
		usr.Connection.WriteJSON(msg)
	}
}
func (cs *chatServer) emitPrivateMessage(targetUserName string, msg message) {
	for _, ch := range cs.channels {
		for _, usr := range ch.Users {
			if usr.Username != targetUserName {
				continue
			}

			var messg struct {
				Type string
				Body interface{}
				From string
			}

			messg.Type = privateType
			messg.Body = msg.Text
			messg.From = msg.Username
			println(usr.Connection)
			usr.Connection.WriteJSON(messg)
			break
		}
	}
}

func (cs *chatServer) emitBroadcast(msg message) {
	for _, usr := range cs.users {
		usr.Connection.WriteJSON(msg)
	}
}

func (cs *chatServer) createChannel(chanName string) *channel {
	newChannel := channel{
		ChannelName: chanName,
		Users:       make([]user, 0),
		Messages:    make([]message, 0),
	}
	return &newChannel
}

func (cs *chatServer) catchUser(previousCh *channel, conn *websocket.Conn) user {
	var usr user
	for i := range previousCh.Users {
		if previousCh.Users[i].Connection != conn {
			continue
		}
		usr = previousCh.Users[i]
		previousCh.Users = append(previousCh.Users[:i], previousCh.Users[i+1:]...)
		break
	}

	return usr
}
func (cs *chatServer) addUsertoChannel(chName string, usr user) {
	for i := range cs.channels {
		if chName != cs.channels[i].ChannelName {
			continue
		}
		cs.channels[i].Users = append(cs.channels[i].Users, usr)
	}
}
func (cs *chatServer) addUserName(usrName string, conn *websocket.Conn) {
	for i := range cs.users {
		if cs.users[i].Connection != conn {
			continue
		}
		cs.users[i].Username = usrName
		break
	}
}
func (cs *chatServer) findUserName(conn *websocket.Conn) string {
	var usrName string
	for _, usr := range cs.users {
		if usr.Connection != conn {
			continue
		}
		usrName = usr.Username
		break
	}
	return usrName
}

func (cs *chatServer) emitPrivateMSG(targetUsr string, fromUser string, messg string) {
	for _, usr := range cs.users {
		if targetUsr != usr.Username {
			continue
		}
		var msg struct {
			Type string
			Body interface{}
			From string
		}
		msg.Type = privateType
		msg.Body = messg
		msg.From = fromUser

		usr.Connection.WriteJSON(msg)
	}

}

func (cs *chatServer) findUserChannelByName(targetUserName string) *channel {
	var targetUserChannel *channel
searchLoop:
	for _, ch := range cs.channels {
		for _, usr := range ch.Users {
			if usr.Username != targetUserName {
				continue
			}
			targetUserChannel = &ch
			break searchLoop
		}
	}
	return targetUserChannel
}

func (cs *chatServer) sendUsersList(conn *websocket.Conn) error {
	var msg struct {
		Type string
		Body interface{}
	}
	msg.Type = userType
	msg.Body = cs.users
	return conn.WriteJSON(msg)
}

func (cs *chatServer) broadcastUsersList() {
	for _, usr := range cs.users {
		err := cs.sendUsersList(usr.Connection)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (cs *chatServer) sendChannelsList(conn *websocket.Conn) error {
	var msg struct {
		Type string
		Body interface{}
	}
	msg.Type = channelType
	msg.Body = cs.channels
	return conn.WriteJSON(msg)
}

func (cs *chatServer) broadcastChannelsList() {
	for _, conn := range cs.users {
		err := cs.sendChannelsList(conn.Connection)
		if err != nil {
			log.Fatal(err)
		}
	}
}
