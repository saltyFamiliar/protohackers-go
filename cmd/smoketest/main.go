package main

import "protohackers-go/internal/smoketest"

func main() {
	go smoketest.TcpEcho()
	go smoketest.UdpEcho()
	select {}
}
