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

type WorkerPool struct {
	tasks chan net.Conn //channel used to send and receive network related operations
}

func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{tasks: make(chan net.Conn, 100)} // use of & to obtain the memory address of the struct
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return pool
}

// worker method since a receiver was added between func and the name function
func (p *WorkerPool) worker() {
	for conn := range p.tasks {
		handleConnection(conn)
	}
}

// func readMessage(conn net.Conn) ([]byte, error) {
//	lenBuff := make([]byte, 4)
//	_, err := io.ReadFull(conn, lenBuff)
//	if err != nil {
//		return nil, err
//	}
//	length := binary.BigEndian.Uint32(lenBuff)
//	data := make([]byte, length)
//	_, err = io.ReadFull(conn, data)
//	if err != nil {
//		fmt.Errorf("Error reading Message from cliente %w", err)
//	}
//	fmt.Printf("Received %v\n", data)
//	return data, err
//}

func handleConnection(conn net.Conn) {
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

func response(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
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

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Unable to connect %v", err)
	}
	defer ln.Close()
	pool := NewWorkerPool(10)
	log.Printf("Listening to TCP connections on port :8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Unable to listen")
			continue
		}
		pool.tasks <- conn
		//go response(conn)
	}

}
