package main

import (
	"log"
	"net"
	. "progohackers/internal/means2anend"
	. "progohackers/internal/utils"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	msgs := make(PriceHistory, 0)

	for i := 0; i < 10; i++ {
		var msg RequestPacket

		_, err := conn.Read(msg[:])
		if err != nil {
			return
		}

		switch requestType := msg.RequestType(); requestType {
		case Insert:
			msgs = append(msgs, msg)
			timeStamp := msg.TimeStamp()
			price := msg.Price()
			log.Printf("Received insertion byte: %x - %c %d %d\n",
				msg,
				requestType,
				timeStamp,
				price,
			)
			msgs = InsertRequestPacket(msgs, msg)
		case Query:
			if len(msgs) == 0 {
				log.Println("Received query request but no packets to query")
				continue
			}

			startTimeRange := msg.StartTime()
			endTimeRange := msg.EndTime()
			log.Printf("Received insertion byte: %x - %c %d %d\n",
				msg,
				requestType,
				startTimeRange,
				endTimeRange,
			)
		default:
			log.Println("Invalid first byte")

		}
	}
	println(msgs)
	conn.Close()
}

func main() {
	l, err := net.Listen("tcp", ":7")
	Must("bind to local address", err)

	for {
		conn, err := l.Accept()
		Must("accept incoming connection", err)
		go handleConnection(conn)
	}
}
