package chatclient

import (
	"fmt"
	"net"
)

type Client struct {
	NickName string
	MsgChan  chan string
}

func NewClient(name string) *Client {
	return &Client{
		NickName: name,
		MsgChan:  make(chan string),
	}
}

func (cl *Client) Writer(conn net.Conn) {
	for msg := range cl.MsgChan {
		fmt.Fprintln(conn, msg)
	}
}

func (cl *Client) ChangeNickName(name string) string {
	oldName := cl.NickName
	cl.NickName = name
	cl.MsgChan <- "Ok. Now you are " + cl.NickName
	return fmt.Sprintf("%s has changed its name to %s", oldName, cl.NickName)
}
