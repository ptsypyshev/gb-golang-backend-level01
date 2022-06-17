package main

import (
	"log"
	"net"

	cs "github.com/ptsypyshev/gb-golang-backend-level01/lesson02/02-chat/internal/chatserver"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	srv := cs.NewServer()
	log.Printf("%s is running", srv.ServerName)
	go srv.Broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go srv.HandleConn(conn)
	}
}
