package main

import (
	"encoding/json"
	"log"
	"net"
)

var (
	userIP map[userID]string
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
		handleLogin(conn, msg.Content[0].([]byte))
	case LOGOUT:
		handleLogout(conn, msg.Content[0].([]byte))
	case REGISTER:
		handleRegister(conn, msg.Content[0].([]byte))
	case CHATMESSAGE:
		handleChatMessage(conn, msg.Content)
	case GROUPMESSAGE:
		handleGroupMessage(conn, msg.Content)
	}
}

func handleLogin(conn net.Conn, content []byte) {
	msg := User{}

	var resMsg string
	var statusCode int

	json.Unmarshal(content, &msg)
	if ok, err := msg.Login(); ok {
		// store userID and IP in map
		userIP[msg.ID] = conn.RemoteAddr().String()
		log.Println("userID: ", msg.ID, "\tIP: ", userIP[msg.ID])

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

	Response(conn, statusCode, resMsg)
}

func handleLogout(conn net.Conn, content []byte) {
	msg := User{}
	json.Unmarshal(content, &msg)
	if _, ok := userIP[msg.ID]; ok {
		delete(userIP, msg.ID)

		log.Println("userID: ", msg.ID, "\tIP: ", userIP[msg.ID])
		Response(conn, LOGOUT, "Logout Success")

		// close connection
		conn.Close()
	} else {
		Response(conn, ERROR, "User Not Found")
	}
}

func handleRegister(conn net.Conn, content []byte) {
	msg := User{}
	json.Unmarshal(content, &msg)
	if msg.CheckExist() {
		Response(conn, ERROR, "User Already Exist")
	} else {
		if err := msg.Add(); err != nil {
			Response(conn, ERROR, "Register Fail")
		} else {
			Response(conn, REGISTER, "Register Success")
		}
	}
}

func handleChatMessage(conn net.Conn, content []interface{}) {
	msg := []Message{}
	for i, v := range content {
		json.Unmarshal(v.([]byte), &msg[i])
		if _, ok := userIP[msg[i].SenderID]; ok {
			// msg[i].Store()
		} else {
			Response(conn, ERROR, "Sender Not Found")
		}
	}
}

func handleGroupMessage(conn net.Conn, content []interface{}) {
	msg := []Message{}
	for i, v := range content {
		json.Unmarshal(v.([]byte), &msg[i])
		if _, ok := userIP[msg[i].SenderID]; ok {
			// msg[i].Store()
		} else {
			Response(conn, ERROR, "Sender Not Found")
		}
	}
}

func Response(conn net.Conn, statusCode int, content string) {
	resContent := []interface{}{[]byte(content)}
	res := clientMessageAPI{Type: statusCode, Content: resContent}
	resJSON, _ := json.Marshal(res)
	conn.Write([]byte(resJSON))
}
