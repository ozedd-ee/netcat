package src

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
)

func Exec(port string, command string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	defer listener.Close()
	fmt.Printf("Listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		fmt.Println("Connection established with", conn.RemoteAddr())

		// Handle the connection in a separate goroutine
		go handleConnection(conn, command)
	}
}

func handleConnection(conn net.Conn, command string) {
	defer conn.Close()

	cmd := exec.Command(command)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Failed to get stdin pipe: %v\n", err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to get stdout pipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start command: %v\n", err)
		return
	}

	// Copy data from the connection to the command's stdin
	go func() {
		defer stdin.Close()
		io.Copy(stdin, conn)
	}()

	// Copy data from the command's stdout to the connection
	go func() {
		io.Copy(conn, stdout)
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v\n", err)
	}
	fmt.Printf("Connection with %s closed.\n", conn.RemoteAddr())
}
