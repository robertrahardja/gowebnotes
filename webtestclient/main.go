package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Write to the server
	fmt.Fprintf(conn, "GET /index.html'\n")

	// Read the response from the server
	bs := make([]byte, 1024)
	n, err := conn.Read(bs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bs[:n]))
}
