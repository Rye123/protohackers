package main

import (
	"net"
	"fmt"
)

const SERVER_PORT = 6969
const DEFAULT_TIMEOUT = 5

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
			buf := make([]byte, 256)
			for {
				count, err := c.Read(buf)
				if err != nil { // Includes io.EOF
					c.Close()
					fmt.Println("Connection closed.")
					break
				}
				if count == 0 {
					continue
				}
				fmt.Printf("Received: %v\n", string(buf))
				
				_, err = c.Write(buf)
				if err != nil {
					fmt.Printf("Error in c.Write: %v\n", err)
				}
				fmt.Printf("Sent: %v\n", string(buf))
			}
		}(conn)
	}
}
