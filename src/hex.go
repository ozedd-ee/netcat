package src

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

func ListenHex(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	defer listener.Close()
	fmt.Println("Listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleHexConnection(conn)
	}
}

func handleHexConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection established with", conn.RemoteAddr())
	
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client.")
				return
			}
			log.Printf("Error reading from connection: %v\n", err)
			return
		}
		_, err = conn.Write([]byte(buffer[:n]))
		if err != nil {
			fmt.Printf("\nError sending message to %v: %v\n", conn.RemoteAddr().String(), err)
		}
		fmt.Println("Received data (hex):")
		fmt.Println(hex.Dump(buffer[:n]))
	}
}
