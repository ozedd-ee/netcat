package src

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func handleConnectionTCP(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection established with", conn.RemoteAddr())

	buffer := make([]byte, 1024)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Messenger stopped for ", conn.RemoteAddr())
				return
			default:
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("Enter message for %v: ", conn.RemoteAddr().String())
				msg, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						return
					}
					fmt.Printf("\nError reading input message: %v\n", err)
					return
				}

				_, err = conn.Write([]byte(msg))
				if err != nil {
					fmt.Printf("\nError sending message to %v: %v\n", conn.RemoteAddr().String(), err)
				}
			}
		}
	}()

	go func() {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				// Check if the connection was closed (EOF is expected when the client disconnects)
				if err == io.EOF {
					fmt.Println("\nClient disconnected:", conn.RemoteAddr())
					close(done)
					return
				}

				fmt.Printf("\nError reading from connection %v: %v\n", conn.RemoteAddr(), err)
				close(done)
				return
			}

			fmt.Printf("\nMessage received from %s: %s\nEnter message for %s: ", conn.RemoteAddr(), string(buffer[:n]), conn.RemoteAddr())

			_, err = conn.Write(buffer[:n])
			if err != nil {
				fmt.Printf("\nError sending response to %s: %v\n", conn.RemoteAddr(), err)
				return
			}
		}
	}()
	<-done
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
