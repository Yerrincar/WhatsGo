package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type WorkerPool struct {
	tasks   chan net.Conn //channel used to send and receive network related operations
	manager *clientManager
}

type clientManager struct {
	clients map[net.Conn]bool
}

func NewWorkerPool(size int, manager *clientManager) *WorkerPool {
	pool := &WorkerPool{tasks: make(chan net.Conn, 100), manager: manager} // use of & to obtain the memory address of the struct
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return pool
}

// worker method since a receiver was added between func and the name function
func (p *WorkerPool) worker() {

	for conn := range p.tasks {
		handleConnection(conn, p.manager)
	}
}

func handleConnection(conn net.Conn, c *clientManager) {
	defer conn.Close()
	c.clients[conn] = true
	c.broadcastMessage(conn, fmt.Sprintf("New connection from %s\n", conn.RemoteAddr().String()))
	for {
		//header of the message first // lenght-prefix protocol
		header := make([]byte, 4)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			log.Printf("Error reading header  %v\n", err)
			return
		}
		length := binary.BigEndian.Uint32(header) //defining endianness
		message := make([]byte, length)
		n, err := io.ReadFull(conn, message)
		if err != nil || n == 0 {
			delete(c.clients, conn)
			c.broadcastMessage(conn, fmt.Sprintf("Client disconnected %s\n", conn.RemoteAddr().String()))
			log.Fatalf("Couldn't read message %v\n", err)
		}
		c.broadcastMessage(conn, fmt.Sprintf("%s: %s", conn.RemoteAddr().String(), string(message[:n])))
	}
}

func (c *clientManager) broadcastMessage(sender net.Conn, message string) {
	for conn := range c.clients {
		if conn == sender {
			continue
		}
		length := len(message)
		header := make([]byte, 4)
		binary.BigEndian.PutUint32(header, uint32(length))
		finalMessage := append(header, message...)
		_, err := conn.Write(finalMessage)
		if err != nil {
			delete(c.clients, conn)
			fmt.Printf("Could not write message")
		}
	}
}

func main() {
	const addr = "localhost:8080"
	fmt.Println("Starting TCP Server on port " + addr)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("Could not resolve address: %s", err.Error())
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}

	defer ln.Close()
	manager := &clientManager{clients: make(map[net.Conn]bool)}
	pool := NewWorkerPool(10, manager)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Fatalf("Unable to listen")
			continue
		}
		pool.tasks <- conn
	}

}
