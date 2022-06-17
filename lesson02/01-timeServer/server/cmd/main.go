package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	go readMsgFromConsole(messages)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			fmt.Println(msg)
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func readMsgFromConsole(msg chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg <- scanner.Text()
	}

	if scanner.Err() != nil {
		fmt.Println("Cannot read from console!")
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	entering <- ch
	input := bufio.NewScanner(conn)
	for input.Scan() {
	}
	leaving <- ch
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for {
		select {
		case msg := <-ch:
			_, err := io.WriteString(conn, msg)
			if err != nil {
				return
			}
		default:
			_, err := io.WriteString(conn, time.Now().Format("15:04:05\n\r"))
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
