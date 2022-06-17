package chatserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	cc "github.com/ptsypyshev/gb-golang-backend-level01/lesson02/02-chat/internal/chatclient"
)

type Server struct {
	ServerName string
	Clients    map[cc.Client]struct{}
	entering   chan cc.Client
	leaving    chan cc.Client
	messages   chan string
}

func NewServer() *Server {
	return &Server{
		ServerName: "SimpleChat v0.1",
		Clients:    make(map[cc.Client]struct{}),
		entering:   make(chan cc.Client),
		leaving:    make(chan cc.Client),
		messages:   make(chan string),
	}
}

func (s *Server) Broadcaster() {
	for {
		select {
		case msg := <-s.messages:
			log.Println(msg)
			for cli := range s.Clients {
				cli.MsgChan <- msg
			}
		case cli := <-s.entering:
			s.Clients[cli] = struct{}{}
		case cli := <-s.leaving:
			delete(s.Clients, cli)
			close(cli.MsgChan)
		}
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	cl := cc.NewClient(conn.RemoteAddr().String())
	go cl.Writer(conn)
	cl.MsgChan <- "You are " + cl.NickName
	cl.MsgChan <- "To set nickname use command /nick [your_nickname], for example /nick Pavel"
	s.messages <- cl.NickName + " has arrived"
	s.entering <- *cl
	input := bufio.NewScanner(conn)
	for input.Scan() {
		text := input.Text()
		if strings.Contains(text, "/nick") {
			newNickName := strings.Replace(text, "/nick ", "", 1)
			s.messages <- cl.ChangeNickName(newNickName)
		} else {
			s.messages <- fmt.Sprintf("%s: %s", cl.NickName, text)
		}
	}
	s.leaving <- *cl
	time.Sleep(time.Second) // Small pause to delete client from broadcast pool
	s.messages <- cl.NickName + " has left"
	conn.Close()
}
