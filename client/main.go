package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalf("Unable to connect to port :8080")
	}
	msg := []byte("Hola\n")
	length := uint32(len(msg))

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, length)
	for i := range 10 {
		conn.Write(lenBuf)
		//os.Stdin.Write(msg)
		//os.Stdin.Write(lenBuf)
		conn.Write(msg)
		fmt.Printf("Enviando mensaje %v\n", i)
		log.Printf("Respuesta del servidor: %v\n", conn)
	}
}
