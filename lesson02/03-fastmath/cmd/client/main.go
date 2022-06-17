package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/ptsypyshev/gb-golang-backend-level01/lesson02/03-fastmath/internal/conf"
)

func main() {
	config := conf.GetConfig()
	conn, err := net.Dial(config["proto"], config["host"]+":"+config["port"])
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
