package main

import "progohackers/internal/smoketest"

func main() {
	go smoketest.TcpEcho()
	go smoketest.UdpEcho()
	select {}
}
