package tests_common

import (
	"fmt"
	"net"
)

func GetFreeLocalAddr() (string, string, int) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port
	return fmt.Sprintf("localhost:%d", port), "localhost", port
}
