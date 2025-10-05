package main

import (
	"net"
	"fmt"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		fmt.Println("Usage: ./main.go <port_num>")
		return
	}

	PORT := ":" + arguments[1]
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("Failure to open port at %s\n", PORT)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {

}
