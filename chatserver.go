package main

import (
	"github.com/gorilla/websocket"
)

var users []user
var channels []channel

type message struct {
	Text string
}

type user struct {
	Connection *websocket.Conn
}

type channel struct {
	ChannelName string
	Users       []user
	Messages    []message
}

func receiveAndSendMessage(msg string, userConn *websocket.Conn) {
	receivedMessage := message{
		Text: msg,
	}
searchLoop:
	for _, ch := range channels {
		for j := range ch.Users {
			if userConn != ch.Users[j].Connection {
				continue
			}
			ch.Messages = append(ch.Messages, receivedMessage)
			emitMessages(ch, receivedMessage)
			break searchLoop
		}

	}
}

func findChannel(channelName string) *channel {
	for _, ch := range channels {
		if ch.ChannelName != channelName {

		}
	}
}

func emitMessages(ch channel, msg message) {
	for _, usr := range ch.Users {
		usr.Connection.WriteJSON(msg)
	}
}

func createAndAddChannel(chanName string, conn *websocket.Conn) {
	newChannel := channel{
		ChannelName: chanName,
		Users:       make([]user, 0),
		Messages:    make([]message, 0),
	}
	channels = append(channels, newChannel)
	UpdateUserChannel(newChannel, conn)
}

func UpdateUserChannel(nextChannel channel, conn *websocket.Conn) {
usrLoop:
	for i, ch := range channels {
		for j := range ch.Users {
			if ch.Users[j].Connection != conn {
				continue
			}
			extractedConn := ch.Users[j]
			channels[i].Users = append(ch.Users[:j], ch.Users[j+1:]...)
			nextChannel.Users = append(nextChannel.Users, extractedConn)
			break usrLoop
		}
	}
}
