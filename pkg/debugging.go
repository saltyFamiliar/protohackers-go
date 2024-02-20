package pkg

import (
	"bufio"
	"net"
)

func printRawBytes(c net.Conn) {
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}(c)

	reader := bufio.NewReader(c)

	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			println(err.Error())
			break
		}
		println(string(msg))
	}
}
