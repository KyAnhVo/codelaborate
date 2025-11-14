package main

import (
	"net"
	"fmt"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// unexpected error 
	defer func() {
		wg.Wait()
	} ()

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

	CreateRoomState()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		
		wg.Add(1)
		remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
		fmt.Printf("Connected with %s:%d\n", remoteAddr.IP, remoteAddr.Port)
		go HandleConnection(&wg, conn)
	}

}


