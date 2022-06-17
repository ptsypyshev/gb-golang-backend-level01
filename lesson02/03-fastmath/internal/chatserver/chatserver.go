// Package chatserver can be used to serve multiple connections of chat clients
package chatserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	fm "github.com/ptsypyshev/gb-golang-backend-level01/lesson02/03-fastmath/internal/fastmath"

	cc "github.com/ptsypyshev/gb-golang-backend-level01/lesson02/03-fastmath/internal/chatclient"
)

const (
	ServerName              = "SimpleChat and MathTasker"
	ServerVersion           = "0.2"
	MathTaskDurationSeconds = 10
)

// Server is a main struct for chatserver package
type Server struct {
	ServerName string
	MathTask   fm.MathTask
	Clients    map[cc.Client]struct{}
	entering   chan cc.Client
	leaving    chan cc.Client
	messages   chan string
}

// NewServer is a constructor for a Server struct
func NewServer() *Server {
	return &Server{
		ServerName: fmt.Sprintf("%s v%s", ServerName, ServerVersion),
		MathTask:   fm.MathTask{},
		Clients:    make(map[cc.Client]struct{}),
		entering:   make(chan cc.Client),
		leaving:    make(chan cc.Client),
		messages:   make(chan string),
	}
}

// Broadcaster is a method that handle clients pool to broadcast some messages to them
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

// HandleConn is a method to handle a single client connection
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
		} else if s.MathTask.GetQuestion() != "" && s.MathTask.GetAnswer() == text {
			s.messages <- fmt.Sprintf("Congratulations! %s answered rightly! %s%s", cl.NickName, s.MathTask.GetQuestion(), s.MathTask.GetAnswer())
			s.MathTask.SetAll("", "") // No one more player can be winner
		} else {
			s.messages <- fmt.Sprintf("%s: %s", cl.NickName, text)
		}
	}
	s.leaving <- *cl
	time.Sleep(time.Second) // Small pause to delete client from broadcast pool
	s.messages <- cl.NickName + " has left"
	conn.Close()
}

// MathTasker is a method that creates math task and sends it to all clients
func (s *Server) MathTasker() {
	for {
		if len(s.Clients) > 0 {
			s.MathTask.Generate()
			s.messages <- s.MathTask.GetQuestion()
		}
		time.Sleep(MathTaskDurationSeconds * time.Second)
	}
}
