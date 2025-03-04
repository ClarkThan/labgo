package main

import (
	"log"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("reply from server"))
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	lnFile, err := ln.(*net.TCPListener).File()
	if err != nil {
		log.Fatalf("file error: %v", err)
	}
	fd := lnFile.Fd()
	log.Println("listen file fd: ", fd)
	if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		panic(err)
	}
	// SO_REUSEPORT
	if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		switch clientConn := conn.(type) {
		case *net.TCPConn:
			f1, err := clientConn.File()
			if err != nil {
				log.Fatalf("file error: %v", err)
			}
			f2, err := clientConn.File()
			log.Printf("TCP conn addr: %v  %d  %d", clientConn.RemoteAddr(), f1.Fd(), f2.Fd())
			// f2.WriteString("hello world ---")
			syscall.Write(int(f1.Fd()), []byte("hello world --- 1"))
			f1.Close()
			syscall.Write(int(f2.Fd()), []byte("hello world --- 2"))
			f2.Close()
		case *net.UnixConn:
			log.Printf("Unix conn addr: %v", clientConn.RemoteAddr())
		default:
			log.Fatalf("unknown connection type: %T", conn)
		}

		if err != nil {
			log.Fatalf("accept error: %v", err)
		}

		go handleConn(conn)
	}
}
