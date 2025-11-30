package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func connectToServer() (net.Conn, error) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal("Unable to connect to port :8080", err)
	}

	return conn, nil
}
func writter(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Print("Impossible to read", err)
		}
		message := []byte(input)
		length := len(message)
		header := make([]byte, 4)
		binary.BigEndian.PutUint32(header, uint32(length))
		finalMessage := append(header, message...)
		conn.Write(finalMessage)
	}
}

func listener(conn net.Conn) {

	for {
		//header of the message first // lenght-prefix protocol
		header := make([]byte, 4)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			log.Printf("Error reading header  %v\n", err)
			if errors.Is(err, io.EOF) {
				fmt.Errorf("Error reading header from the server: %w", err)
			}
		}
		length := binary.BigEndian.Uint32(header) //defining endianness
		message := make([]byte, length)
		_, err = io.ReadFull(conn, message)
		if err != nil {
			log.Fatalf("Couldn't read message %v\n", err)
		}
	}
}

func main() {
	conn, err := connectToServer()
	if err != nil {
		log.Fatalf("Unable to connect to port :8080")
	}
	go writter(conn)
	go listener(conn)
	select {}
}
