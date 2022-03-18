package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"net"
	"xorm.io/xorm"
)

var (
	engine *xorm.Engine
	userIP map[string]string
)

const (
	LOGIN = iota
	LOGOUT
	REGISTER
	CHATMESSAGE
	GROUPMESSAGE
)

type clientMessageAPI struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

func main() {
	go func() {
		engine = initDB()
	}()
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

func handleClient(conn net.Conn) {
	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			log.Fatal(err)
		} else {
			if n == 0 {
				continue
			} else {
				go handleMessage(conn, buf[0:n])
			}
		}
		// _, err2 := conn.Write(buf[0:n])
		// if err2 != nil {
		// 	return
		// }
	}
}

func handleMessage(conn net.Conn, buf []byte) {
	var msg clientMessageAPI
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		log.Fatal(err)
	}
	switch msg.Type {
	case LOGIN:
		handleLogin(conn, msg.Content)
	case LOGOUT:
		handleLogout(conn, msg.Content)
	case REGISTER:
		handleRegister(conn, msg.Content)
	case CHATMESSAGE:
		handleChatMessage(conn, msg.Content)
	case GROUPMESSAGE:
		handleGroupMessage(conn, msg.Content)
	}
}
