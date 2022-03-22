package main

import (
	"log"
	"net"
)

const (
	LOGIN = iota
	LOGOUT
	REGISTER
	CHATMESSAGE
	GROUPMESSAGE
	ERROR
)

type clientMessageAPI struct {
	Type    int           `json:"type"`
	Content []interface{} `json:"content"`
}

func main() {
	go InitDB()

	// socket server
	service := ":1201"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleClient(conn)
	}
}
