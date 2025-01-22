package src

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type ClientSession struct {
	Address    net.Addr
	LastActive time.Time
}

const (
	timeoutPeriod   = 30 * time.Second
	cleanupInterval = 10 * time.Second
)

var (
	sessionMap = make(map[string]*ClientSession)
	mutex      sync.Mutex
)

func handleConnectionUDP(conn net.PacketConn) {
	defer conn.Close()

	go cleanupInactiveSessions()

	go handleServerInput(conn)

	buffer := make([]byte, 1024)

	for {
		// Read incoming packet
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("Error reading from client: %v\n", err)
			continue
		}

		// Update session information
		updateSession(addr)

		go handleClientMessage(conn, addr, buffer[:n])
	}
}

func handleClientMessage(conn net.PacketConn, addr net.Addr, data []byte) {
	fmt.Printf("Message received from %s: %s\n", addr.String(), data)

	response := fmt.Sprintf("You sent: %s", data)
	_, err := conn.WriteTo([]byte(response), addr)
	if err != nil {
		fmt.Printf("Error sending response to %s: %v\n", addr.String(), err)
	}
}

func updateSession(addr net.Addr) {
	mutex.Lock()
	defer mutex.Unlock()

	clientKey := addr.String()
	session, exists := sessionMap[clientKey]
	if !exists {
		// Create a new session if the client is new
		sessionMap[clientKey] = &ClientSession{
			Address:    addr,
			LastActive: time.Now(),
		}
		fmt.Printf("New client connected: %s\n", clientKey)
	} else {
		// Update the last active time for an existing client
		session.LastActive = time.Now()
	}
}

func cleanupInactiveSessions() {
	for {
		time.Sleep(cleanupInterval)

		mutex.Lock()
		for key, session := range sessionMap {
			if time.Since(session.LastActive) > timeoutPeriod {
				// Remove inactive client
				delete(sessionMap, key)
				fmt.Printf("\nRemoved inactive client: %s\n", key)
			}
		}
		mutex.Unlock()
	}
}

func handleServerInput(conn net.PacketConn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to broadcast: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading server input: %v\n", err)
			continue
		}
		// Broadcast the message to all connected clients
		broadcastMessage(conn, input)
	}
}

func broadcastMessage(conn net.PacketConn, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, session := range sessionMap {
		_, err := conn.WriteTo([]byte(message), session.Address)
		if err != nil {
			fmt.Printf("Error broadcasting to %s: %v\n", session.Address, err)
		} else {
			fmt.Printf("Broadcasted to %s: %s\n", session.Address, message)
		}
	}
}
