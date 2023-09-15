package main

import (
	"net"
	"io"
	"fmt"
	"time"
)

const SERVER_PORT = 6969
const DEFAULT_TIMEOUT = 5

// Receives until '\n' through the given connection.
func recvMsg(conn net.Conn) (msg string, err error) {
	message := make([]byte, 0, 1024)
	
	for {
		buf := make([]byte, 256)
		_, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				return string(message), nil
			}
			return "", err
		}

		for _, c := range buf {
			message = append(message, c)
			if c == '\n' {
				return string(message), nil
			}
		}
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	}
}

func main() {
	fmt.Printf("Running on %d.\n", SERVER_PORT)
	listen_sock, err := net.Listen("tcp", fmt.Sprintf(":%d", SERVER_PORT))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer listen_sock.Close()

	for {
		conn, err := listen_sock.Accept()
		fmt.Printf("Received connection from %v.\n", conn.RemoteAddr())
		if err != nil {
			fmt.Printf("Connection Error: %v\n", err)
			continue
		}

		go func(c net.Conn) {
			msg, err := recvMsg(c)
			if err != nil {
				fmt.Printf("Connection Receive Error: %v\n", err)
				return
			}

			fmt.Printf("Received message: %v\n", msg)
			_, err = c.Write([]byte(msg))
			if err != nil {
				fmt.Printf("Send Error: %v\n", err)
			}
			fmt.Printf("Sent: %v\n", msg)
			
			c.Close()
		}(conn)
	}
}
