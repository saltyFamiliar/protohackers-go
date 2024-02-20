package primetime

import (
	"encoding/json"
	"net"
	"testing"
)

func testNumber(t *testing.T, n int, want bool) {
	conn, err := net.Dial("tcp", ":7")
	if err != nil {
		t.Fatal(err)
	}

	floatn := float64(n)

	req := primeReq{
		Method: "isPrime",
		Number: &floatn,
	}

	err = json.NewEncoder(conn).Encode(req)
	if err != nil {
		t.Fatal(err)
	}

	var resp primeResp
	err = json.NewDecoder(conn).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Method != "isPrime" {
		t.Errorf("got method %q, want %q", resp.Method, "isPrime")
	}

	if resp.IsPrime != want {
		t.Errorf("got isPrime=%t for n=%d, want %t", resp.IsPrime, n, want)
	}

	conn.Close()
}

func TestTcpPrimeTest(t *testing.T) {

	// Start the server in a goroutine
	go TcpPrimeTest()

	// Open some connections to test it
	for _, n := range []int{2, 3, 5, 7, 11, 13} {
		testNumber(t, n, true)
	}

	// Check some non-prime numbers
	for _, n := range []int{4, 6, 8, 9} {
		testNumber(t, n, false)
	}
}
