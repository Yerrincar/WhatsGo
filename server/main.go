package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
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

func readMessage(conn net.Conn) ([]byte, error) {
	lenBuff := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuff)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuff)
	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	fmt.Printf("Received %v\n", data)
	return data, err
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		conn.SetReadDeadline(time.Now().Add(20 * time.Second)) //Timeout if any read operation surpass 6 seconds
		data, err := readMessage(conn)
		if err != nil {
			log.Print("Wait time over: ", err)
			return // close client connection
		}
		fmt.Printf("Received %v\n", data)
		conn.Write([]byte("Message received\n"))
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
	}

}
