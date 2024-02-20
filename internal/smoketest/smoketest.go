package smoketest

import (
	"io"
	"net"
)

func TcpEcho() {
	l, err := net.Listen("tcp", ":7")
	if err != nil {
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			io.Copy(c, c)
			c.Close()
		}()
	}
}

func UdpEcho() {
	conn, err := net.ListenPacket("udp", ":7")
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		_, err = conn.WriteTo(buf[:n], addr)
		if err != nil {
			panic(err)
		}
	}
}
