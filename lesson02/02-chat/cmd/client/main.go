package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	_, errIn := io.Copy(conn, os.Stdin) // until you send ^Z
	if errIn != nil {
		fmt.Println(errIn)
		return
	}
	fmt.Printf("%s: exit", conn.LocalAddr())
}
