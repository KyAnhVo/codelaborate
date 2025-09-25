package main

import (
	"fmt"
	"net/http"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go serverTCP(&wg)
	go serverHTTP(&wg)
	wg.Wait()
}

func serverTCP(wg *sync.WaitGroup) {
	defer wg.Done()
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Cannot listen on port 80")
		return
	}
	fmt.Println("Started listening")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Cannot accept call")
			conn.Close()
			break
		}
		go sendTCP(conn)
	}
}

func sendTCP(conn net.Conn) {
	fmt.Fprintln(conn, "{\"msg\": \"Hello World\"}")
	conn.Close()
}

func serverHTTP(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc("/", sendWebapp)
	fmt.Println("HTTP server listen")
	err:= http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP: ListendAndServe: ", err)
		return
	}
}

func sendWebapp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP: client at ", r.RemoteAddr)
	filePath := "./../client/index.html"
	http.ServeFile(w, r, filePath)
}
