package dbproxy

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	PROXY_ADDR = "127.0.0.1:3307"
	DB_ADDR    = "127.0.0.1:3306"

	COM_QUERY = byte(0x03)
)

// https://www.youtube.com/watch?v=DU7_MQmRDUs
func Main() {
	proxy, err := net.Listen("tcp", PROXY_ADDR)
	if err != nil {
		fmt.Printf("Listen proxy error: %v\n", err)
		os.Exit(1)
	}

	for {
		conn, err := proxy.Accept()
		if err == nil {
			fmt.Printf("accept error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("new connection from %s\n", conn.RemoteAddr())
		go transport(conn)
	}
}

func transport(conn net.Conn) {
	defer conn.Close()

	dbConn, err := net.Dial("tcp", DB_ADDR)
	if err != nil {
		fmt.Printf("connect to db error: %v\n", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	// 1. forward client request to db
	go intercept(conn, dbConn)

	// 2. forward db response to client blocking
	if _, err := io.Copy(conn, dbConn); err != nil {
		fmt.Printf("copy error: %v\n", err)
	}
}

func intercept(client, db net.Conn) {
	buf := make([]byte, 4096)

	for {
		n, _ := client.Read(buf)
		// 3-size, 1-seq, 1-com code
		// 3 - length of body, 1 - packet sequence number, 1 - command code, etc.
		if n > 5 {
			switch buf[4] {
			case COM_QUERY:
				query := string(buf[5:n])
				newQuery := strings.ReplaceAll(query, "from demo1", "from demo")
				fmt.Printf("orginal query: %s\nto server: %s\n", query, newQuery)
				copy(buf[5:n], []byte(newQuery))
			}
		}
		db.Write(buf[:n])
	}
}
