package src

import (
	"context"
	"fmt"
	"log"
	"net"
)

func TCPListen(network, port string) {
	address := ":" + port
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnectionTCP(conn)
	}
}

func UDPListen(ctx context.Context, network, port string) {
	address := ":" + port
	conn, err := net.ListenPacket(network, address)
	if err != nil {
		fmt.Printf("failed to start server: %v", err)
	}
	fmt.Printf("Server is listening on port %s\n", port)
	handleConnectionUDP(conn)

	<-ctx.Done()
	// Close the connection to unblock `conn.ReadFrom`
	conn.Close()
}
