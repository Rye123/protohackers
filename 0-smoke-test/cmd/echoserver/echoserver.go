package main

import (
	"fmt"
	"net"
	"io"
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
			buf := make([]byte, 1024)
			for {
				count, err := c.Read(buf)
				if count == 0 {
					if err == io.EOF {
						c.Close()
						fmt.Println("Connection closed.")
						break
					}
					if err != nil {
						c.Close()
						fmt.Printf("Connection closed due to error: %v\n", err)
						break
					}
					// Otherwise, continue reading
					continue
				}

				if err != nil && err != io.EOF {
					c.Close()
					fmt.Printf("Connection closed prematurely due to error: %v\n", err)
					break
				}

				buf = buf[:count] // Truncate to exactly number of received bytes

				fmt.Printf("Recv: %v\n", string(buf))

				// Write exactly that number of bytes as output
				_, err = c.Write(buf[:count])
				if err != nil {
					fmt.Printf("Error in c.Write: %v\n", err)
				}
				fmt.Printf("Sent: %v\n", string(buf))
			}
		}(conn)
	}
}
