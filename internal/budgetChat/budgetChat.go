package budgetChat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"unicode"
)

type ChatMessage struct {
	Text       string
	SenderName string
	IsSpecial  bool
}

type SystemMessageCode byte

const (
	DeleteUser = iota
	CreateUser
)

type SystemMessage struct {
	Code   SystemMessageCode
	Sender *User
}

type User struct {
	Name   string
	ChatCh chan string
}

func validateName(name string) bool {
	for _, c := range []rune(name) {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}

func HandleUserSession(conn net.Conn, systemCh chan SystemMessage, globalMsgCh chan ChatMessage) (err error) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	var name string
	if _, err := rw.Write([]byte("What's your name?\n")); err != nil {
		return err
	}
	if err = rw.Flush(); err != nil {
		return err
	}

	if name, err = rw.ReadString('\n'); err != nil {
		return err
	}
	name = name[:len(name)-1]
	if !validateName(name) {
		return fmt.Errorf("invalid name")
	}

	thisUser := User{Name: name, ChatCh: make(chan string, 10)}
	systemCh <- SystemMessage{Code: CreateUser, Sender: &thisUser}
	defer func() {
		systemCh <- SystemMessage{Code: DeleteUser, Sender: &thisUser}
	}()

	// write messages from chatCh to user
	go func() {
		for msg := range thisUser.ChatCh {
			if _, err := rw.Write([]byte(msg)); err != nil {

			}
			if err = rw.Flush(); err != nil {
				log.Println("Unable to flush buffer")
			}
		}
	}()

	// read messages from user to global chat
	for {
		var message string
		if message, err = rw.ReadString('\n'); err != nil {
			globalMsgCh <- ChatMessage{
				Text:       "has left the room\n",
				SenderName: thisUser.Name,
				IsSpecial:  true,
			}
			return err
		}
		log.Printf("{%s}: %s", thisUser.Name, message)

		go func() {
			globalMsgCh <- ChatMessage{Text: message, SenderName: thisUser.Name}
		}()
	}
}
