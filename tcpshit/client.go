package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if errors.Is(err, io.EOF) {
			log.Println("read EOF")
			return
		}

		fmt.Println("receive: ", string(buf[:n]))
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}

	handleClient(conn)
}
