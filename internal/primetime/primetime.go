package primetime

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net"
	"strconv"
)

type primeReq struct {
	Method string   `json:"method"`
	Number *float64 `json:"number,omitempty"`
}

type primeResp struct {
	Method  string `json:"method"`
	IsPrime bool   `json:"prime"`
}

func isPrime(n float64) bool {
	if math.Floor(n) != n {
		return false
	}

	bigN := new(big.Int)
	bigN.SetString(strconv.FormatFloat(n, 'f', 0, 64), 10)
	return bigN.ProbablyPrime(20)
}

func processReq(line []byte) (string, error) {
	if !json.Valid(line) {
		return "invalid json\n", fmt.Errorf("invalid json")
	}

	var req primeReq
	dec := json.NewDecoder(bytes.NewReader(line))
	if err := dec.Decode(&req); err != nil {
		return "error decoding\n", fmt.Errorf("error decoding")
	}
	if dec.More() {
		return "json objects not newline terminated\n", fmt.Errorf("json objects not newline terminated")
	}

	if req.Method != "isPrime" || req.Number == nil {
		return "method or number wrong\n", fmt.Errorf("method or number wrong")
	}

	resp := primeResp{
		Method:  "isPrime",
		IsPrime: isPrime(*req.Number),
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return "error marshalling\n", fmt.Errorf("error marshalling")
	}

	return string(respBytes) + "\n", nil
}

func handleConn(c net.Conn) {
	defer c.Close()

	scanner := bufio.NewScanner(c)

	for scanner.Scan() {
		line := scanner.Bytes()
		resp, err := processReq(line)
		c.Write([]byte(resp))
		fmt.Println(string(line), ": ", resp)
		if err != nil {
			return
		}
	}
}

func TcpPrimeTest() {
	l, err := net.Listen("tcp", ":7")
	if err != nil {
		panic(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		println("conn accepted")
		go handleConn(c)
		//printRawBytes(c)
	}
}
