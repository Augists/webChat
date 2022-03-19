package main

import (
	"encoding/json"
	"log"
	"net"
)

var (
	userIP map[userID]string
)

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
		// handleLogout(conn, msg.Content)
	case REGISTER:
		// handleRegister(conn, msg.Content)
	case CHATMESSAGE:
		// handleChatMessage(conn, msg.Content)
	case GROUPMESSAGE:
		// handleGroupMessage(conn, msg.Content)
	}
}

func handleLogin(conn net.Conn, content []byte) {
	msg := LoginContent{}

	var resMsg string
	var statusCode int

	json.Unmarshal(content, &msg)
	if ok, err := msg.UserLogin(); ok {
		// store userID and IP in map
		userIP[msg.UserID] = conn.RemoteAddr().String()
		log.Println("userID: ", msg.UserID, "\tIP: ", userIP[msg.UserID])

		resMsg = "Login Success"
		statusCode = LOGIN
	} else if err != nil {
		resMsg = "User Not Found"
		statusCode = ERROR
		log.Fatal(err)
	} else {
		resMsg = "Password Not Match"
		statusCode = ERROR
		log.Fatal(err)
	}

	res := clientMessageAPI{Type: statusCode, Content: []byte(resMsg)}
	resJSON, _ := json.Marshal(res)
	conn.Write([]byte(resJSON))
}
