package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const bufferSize = 256

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, bufferSize) // создаем буфер
	for {
		buf = bufClear(buf)
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		_, err = fmt.Fprintln(os.Stdout, strings.TrimSpace(string(buf)))
		if err != nil {
			fmt.Println("cannot write to stdout")
			break
		}
	}
}

// bufClear used to clear buffer content for correct console output
func bufClear(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = ' '
	}
	return b
}
