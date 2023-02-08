package main

import (
	"fmt"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go readMultipleCommands(conn)
	}

}

func readMultipleCommands(conn net.Conn) {

	defer conn.Close()

	for {

		// copied from solutions
		buf := make([]byte, 1024)
		len, err := conn.Read(buf)

		if err != nil {
			fmt.Printf("Error reading: %#v\n", err)
			return
		}

		if len == 0 {
			fmt.Println("Connection Closed")
			return
		}

		conn.Write([]byte("+PONG\r\n"))
	}

}
