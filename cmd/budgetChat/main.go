package main

import (
	"fmt"
	"log"
	"net"
	. "protohackers-go/internal/budgetChat"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", ":7")
	if err != nil {
		panic(err)
	}

	userMap := make(map[string]*User)
	var userMapMutex sync.RWMutex
	globalMsgCh := make(chan ChatMessage, 10)
	go func() {
		for msg := range globalMsgCh {
			var msgView string
			if msg.IsSpecial {
				msgView = fmt.Sprintf("* %s %s", msg.SenderName, msg.Text)
			} else {
				msgView = fmt.Sprintf("[%s] %s", msg.SenderName, msg.Text)
			}
			userMapMutex.RLock()
			for name, user := range userMap {
				if name != msg.SenderName {
					user.ChatCh <- msgView
				}
			}
			userMapMutex.RUnlock()
		}
	}()

	systemCh := make(chan SystemMessage, 10)
	go func() {
		for systemMessage := range systemCh {
			switch systemMessage.Code {
			case DeleteUser:
				userMapMutex.Lock()
				delete(userMap, systemMessage.Sender.Name)
				userMapMutex.Unlock()
				close(systemMessage.Sender.ChatCh)

			case CreateUser:
				userMapMutex.Lock()
				userMap[systemMessage.Sender.Name] = systemMessage.Sender
				userMapMutex.Unlock()

				log.Println("New user joined: ", systemMessage.Sender.Name)
				userListMsg := "* The room contains:  "
				userMapMutex.RLock()
				for name := range userMap {
					if name != systemMessage.Sender.Name {
						userListMsg += name + ", "
					}
				}
				userMapMutex.RUnlock()

				userListMsg = userListMsg[:len(userListMsg)-2] + "\n"
				systemMessage.Sender.ChatCh <- userListMsg

				globalMsgCh <- ChatMessage{
					Text:       fmt.Sprintf("has joined the room\n"),
					SenderName: systemMessage.Sender.Name,
					IsSpecial:  true}
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			err := HandleUserSession(conn, systemCh, globalMsgCh)
			if err != nil {
				log.Println(err)
			}
		}()
	}

}
