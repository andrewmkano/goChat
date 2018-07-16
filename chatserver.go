package main

import (
	"github.com/gorilla/websocket"
)

type chatServer struct {
	users    []user
	channels []channel
}

type message struct {
	Text string
}

type user struct {
	UserName   string
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

func (cs *chatServer) obtainUserFromCh(currentCh string, conn *websocket.Conn) user {
	var usrSwitching user
channelLoop:
	for i := range cs.channels {
		if cs.channels[i].ChannelName != currentCh {
			continue
		}
		channelUsers := cs.channels[i].Users
		for k := range channelUsers {
			if channelUsers[k].Connection != conn {
				continue
			}
			usrSwitching = channelUsers[k]
			cs.channels[i].Users = append(channelUsers[:k], channelUsers[k+1:]...)
			break channelLoop
		}
	}

	return usrSwitching
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

func (cs *chatServer) UpdateUserCh(nextCh string, conn *websocket.Conn) {
	currentChannel := cs.findUserChannel(conn)
	userConn := cs.obtainUserFromCh(currentChannel.ChannelName, conn)
	nextChannel := cs.findChannel(nextCh)
	nextChannel.Users = append(nextChannel.Users, userConn)

}
