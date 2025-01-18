package src

import (
	"fmt"
	"io"
	"net"
)

func handleConnectionTCP(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	fmt.Println("Connection established with", conn.RemoteAddr())

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// Check if the connection was closed (EOF is expected when the client disconnects)
			if err == io.EOF {
				fmt.Println("Client disconnected:", conn.RemoteAddr())
				return
			}
			fmt.Printf("Error reading from connection %v: %v\n", conn.RemoteAddr(), err)
			return
		}

		fmt.Printf("Message received from %s: %s\n", conn.RemoteAddr(), string(buffer[:n]))

		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error sending response to %s: %v\n", conn.RemoteAddr(), err)
			return
		}
	}
}

func handleConnectionUDP(conn net.PacketConn) {
	buffer := make([]byte, 1024)

	for {
		fmt.Println("Waiting for a message...")
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			// Check if the connection was closed; this is expected during shutdown
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				// exit loop
				return
			}

			fmt.Printf("Error reading from connection: %v\n", err)
			continue
		}

		fmt.Printf("Message received from %s: %s\n", addr.String(), string(buffer[:n]))

		_, err = conn.WriteTo([]byte(buffer[:n]), addr)
		if err != nil {
			fmt.Printf("Error sending response: %v\n", err)
		}
	}
}
