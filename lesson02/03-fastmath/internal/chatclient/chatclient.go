// Package chatclient can be used as simple chat client
package chatclient

import (
	"fmt"
	"net"
)

// Client is a main struct for chatclient package
type Client struct {
	NickName string
	MsgChan  chan string
}

// NewClient is a constructor for a Client struct
func NewClient(name string) *Client {
	return &Client{
		NickName: name,
		MsgChan:  make(chan string),
	}
}

// Writer is a method that sends messages to clients
func (cl *Client) Writer(conn net.Conn) {
	for msg := range cl.MsgChan {
		fmt.Fprintln(conn, msg)
	}
}

// ChangeNickName is a method that changes NickName field and returns a string with message
func (cl *Client) ChangeNickName(name string) string {
	oldName := cl.NickName
	cl.NickName = name
	cl.MsgChan <- "Ok. Now you are " + cl.NickName
	return fmt.Sprintf("%s has changed its name to %s", oldName, cl.NickName)
}
