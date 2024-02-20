package main

import (
	"log"
	"math/rand"
	"net"
	. "progohackers/internal/means2anend"
	"time"
)

func sendInsertionPacket(conn net.Conn) error {
	msg := NewInsertionPacket()

	msg.SetTimeStamp(int32(rand.Int63()))
	msg.SetPrice(int32(rand.Int63()))

	msg.SendOverConnection(conn)

	log.Printf("Sent: %x - %c %d %d\n", msg, msg[0], msg.TimeStamp(), msg.Price())

	return nil
}

func sendQueryPacket(conn net.Conn) {
	msg := NewInsertionPacket()

	msg.SetStartTime(int32(rand.Int63()))
	msg.SetEndTime(int32(rand.Int63()))

	msg.SendOverConnection(conn)

	log.Printf("Sent: %x - %c %d %d\n", msg, msg[0], msg.StartTime(), msg.EndTime())
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

		time.Sleep(1 * time.Second)
	}
}
