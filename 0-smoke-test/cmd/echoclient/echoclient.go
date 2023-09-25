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
func recvMsg(conn net.Conn) (msg []byte, err error) {
	message := make([]byte, 0, 1024)
	
	for {
		buf := make([]byte, 256)
		_, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				return message, nil
			}
			return []byte(""), err
		}
		message = append(message, buf[:]...)
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
	recvbuf := make([]byte, 256)
	count := 0
	for count == 0 {
		count, err = conn.Read(recvbuf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection Ended.")
				return
			}
			printErr(fmt.Sprintf("Receive Error: %v", err))
			return
		}
	}
	fmt.Printf("Received: %v\n", string(recvbuf))
	conn.Close()
}
