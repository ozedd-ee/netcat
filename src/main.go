package src

import (
	"fmt"
	"net"
	"time"
)

func Ping(host string, port string) {
	timeout := 2 * time.Second

	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Printf("No TCP server is listening on %s:%s\n", host, port)
		// TCP connection attempt failed, try UDP
		conn, err := net.DialTimeout("udp", address, timeout)
		if err != nil {
			fmt.Printf("No UDP server is listening on %s:%s\n", host, port)
			return 
		}
		// UDP connection attempt successful
		fmt.Printf("UDP Server is listening on %s:%s\n", host, port)
		conn.Close()
		return
	}
	// TCP connection successful
	fmt.Printf("TCP Server is listening on %s:%s\n", host, port)
	conn.Close()
}

func TcpConnect(host string, port string) {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Error connecting to server at %s:%e", address, err)
		return
	}
	handleConnectionTCP(conn)
}

func UdpConnect(host string, port string) {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Error connecting to server at %s:%e", address, err)
		return
	}
	// Type assert to *net.UDPConn since it also implements net.PacketConn
	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		fmt.Println("Failed to type assert to *net.UDPConn")
		return
	}
	handleConnectionUDP(udpConn)
}
