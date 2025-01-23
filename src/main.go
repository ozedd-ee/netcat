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
