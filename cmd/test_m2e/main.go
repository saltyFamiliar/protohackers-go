package main

import (
	"log"
	"math/rand"
	"net"
	. "protohackers-go/internal/means2anend"
	"time"
)

func sendInsertionPacket(conn net.Conn) error {
	msg := NewInsertionPacket(int32(rand.Int63())%1000, int32(rand.Int63())%1000)

	err := msg.SendOverConnection(conn)
	if err == nil {
		log.Printf("Sent: %x - %c %d %d\n", msg, msg[0],
			msg.TimeStamp(), msg.Price())
	} else {
		log.Printf("Failed to send insertion packet")
	}

	return err
}

func sendQueryPacket(conn net.Conn) error {
	msg := NewQueryPacket(int32(rand.Int63())%1000, int32(rand.Int63())%1000)

	err := msg.SendOverConnection(conn)
	if err == nil {
		log.Printf("Sent: %x - %c %d %d\n", msg, msg[0],
			msg.StartTime(), msg.EndTime())
	} else {
		log.Printf("Failed to send query packet")
	}

	return err
}

func main() {
	conn, err := net.Dial("tcp", "localhost:7")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		err := sendInsertionPacket(conn)
		if err != nil {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	for i := 0; i < 10; i++ {
		err := sendQueryPacket(conn)
		if err != nil {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
}
