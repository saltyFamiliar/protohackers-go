package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	. "protohackers-go/internal/means2anend"
	. "protohackers-go/internal/utils"
	"time"
)

func readRequest(conn net.Conn) (msg RequestPacket, err error) {
	for totalBytes, newBytes := 0, 0; totalBytes < 9; totalBytes += newBytes {
		newBytes, err = conn.Read(msg[totalBytes:])
		if err != nil {
			return msg, fmt.Errorf("error reading bytes: %v", err)
		}
	}

	return msg, nil
}

func handleConnection(conn net.Conn) (err error) {
	defer conn.Close()

	msgs := make(PriceHistory, 0)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(time.Second * 15)); err != nil {
			return fmt.Errorf("unable to set deadline")
		}

		var msg RequestPacket
		if msg, err = readRequest(conn); err != nil {
			return err
		}

		switch requestType := msg.RequestType(); requestType {
		case Insert:
			log.Printf("Received insertion packet: %v", msg)
			if msgs, err = InsertRequestPacket(msgs, msg); err != nil {
				return fmt.Errorf("duplicate timestamp")
			}

		case Query:
			if len(msgs) == 0 {
				log.Println("Received query request but no packets to query")
				continue
			}
			log.Printf("Received query packet: %v", msg)

			response := msgs.InRange(msg.StartTime(), msg.EndTime())
			var respBuf [4]byte
			binary.BigEndian.PutUint32(respBuf[:], uint32(response.MeanPrice()))
			if sent, err := conn.Write(respBuf[:]); err != nil || sent != 4 {
				return fmt.Errorf("error sending response")
			}

		default:
			return fmt.Errorf("invalid first byte: %v", msg)
		}
	}
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
