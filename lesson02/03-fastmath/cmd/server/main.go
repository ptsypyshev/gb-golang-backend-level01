package main

import (
	"log"
	"net"

	"github.com/ptsypyshev/gb-golang-backend-level01/lesson02/03-fastmath/internal/conf"

	cs "github.com/ptsypyshev/gb-golang-backend-level01/lesson02/03-fastmath/internal/chatserver"
)

func main() {
	config := conf.GetConfig()
	listener, err := net.Listen(config["proto"], config["host"]+":"+config["port"])
	if err != nil {
		log.Fatal(err)
	}
	srv := cs.NewServer()
	log.Printf("%s is running", srv.ServerName)
	go srv.Broadcaster()
	go srv.MathTasker()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go srv.HandleConn(conn)
	}
}
