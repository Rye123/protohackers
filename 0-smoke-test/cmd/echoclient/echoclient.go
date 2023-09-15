package main

import (
	"fmt"
	"net"
	"io"
	"os"
	"time"
	"bufio"
)

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
			if c == '\n' {
				return string(message), nil
			}
			message = append(message, c)
		}
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	}
}

func printErr(msg string) {
	fmt.Fprintf(os.Stderr, msg + "\n")
}

func main() {
	if len(os.Args) != 3 {
		printErr("Usage: echoclient [host] [port]")
		return
	}
	
	addr := fmt.Sprintf("%v:%v", os.Args[1], os.Args[2])
	
	fmt.Println("Connecting...")
	conn, err := net.DialTimeout("tcp", addr, DEFAULT_TIMEOUT * time.Second)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected to %v.\n", addr)

	// Get user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input: ")
	inp, _ := reader.ReadString('\n')
	
	_, err = conn.Write([]byte(inp))
	fmt.Printf("\nSent: %v\n", inp)
	if err != nil {
		printErr(fmt.Sprintf("Send Error: %v", err))
		return
	}
	
	// Get response
	resp, err := recvMsg(conn)
	if err != nil {
		printErr(fmt.Sprintf("Receive Error: %v", err))
		return
	}

	fmt.Printf("Received: %v\n", resp)


}
