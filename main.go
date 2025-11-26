package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1028)
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //Timeout if any read operation surpass 6 seconds
		n, err := conn.Read(buffer)
		if err != nil {
			log.Print("Wait time over: ", err)
			return // close client connection
		}
		if n == 0 {
			log.Print("No message was sent")
		}
		fmt.Printf("Received %v\n", buffer[:n])
		conn.Write([]byte("Message received\n"))
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Unable to connect %v", err)
	}
	defer ln.Close()
	log.Printf("Listening to TCP connections on port :8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Unable to listen")
			continue
		}
		go handleConnection(conn)
	}

}
