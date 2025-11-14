package main

import (
	"net"
	"fmt"
	"os"
	"sync"
	log "github.com/sirupsen/logrus"
)

func init() {
    // Setup logging
    log.SetFormatter(&log.TextFormatter{
        ForceColors:     true,
        FullTimestamp:   true,
        TimestampFormat: "15:04:05",
    })
    log.SetOutput(os.Stdout)
    log.SetLevel(log.DebugLevel)
}


var wg sync.WaitGroup

func main() {
	// unexpected error 
	defer func() {
		wg.Wait()
	} ()

	arguments := os.Args
	if len(arguments) != 2 {
		log.Error("Usage: ./main.go <port_num>")
		return
	}

	PORT := ":" + arguments[1]
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Errorf("Failure to open port at %s\n", PORT)
		return
	}
	defer listener.Close()

	CreateRoomState()
	for {
		log.Infof("Awaiting conn")
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		
		wg.Add(1)
		remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
		fmt.Printf("Connected with %s:%d\n", remoteAddr.IP, remoteAddr.Port)
		go HandleConnection(&wg, conn)
	}

}


