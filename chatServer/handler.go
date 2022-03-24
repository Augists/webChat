package main

import (
	"encoding/json"
	"log"
	"net"
)

var (
	userIP map[string]string
)

/*
 * There may be some simple method for message type
 * Create 6 ports for those message type
 * and send message to the ports matching message type
 */
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
				/*
				 * Jsonify in python will use <Response>
				 * but json.dumps will use <string> with " surrounded
				 * so we won't deliver " in buf to handleMessage
				 */
				go handleMessage(conn, buf[1:n-1])
				// go handleMessage(conn, buf[0:n])
			}
		}
	}
}

func handleMessage(conn net.Conn, buf []byte) {
	/*
	 * For testing json encoded result in golang
	 */
	// testUser := User{ID: 123, Password: "1234"}
	// testMsg := clientMessageAPI{Type: 0, Content: []interface{}{testUser}}
	// test, _ := json.Marshal(testMsg)
	// fmt.Println(string(test))

	/*
	 * replace ' with " in buf
	 */
	for i, v := range buf {
		if v == '\'' {
			buf[i] = '"'
		}
	}
	// fmt.Println(string(buf))
	var msg clientMessageAPI
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		log.Fatal(err)
	}
	// log.Print(msg.Content)
	// log.Print(msg.Content[0])
	switch msg.Type {
	case LOGIN:
		handleLogin(conn, msg.Content[0].(map[string]interface{}))
	case LOGOUT:
		handleLogout(conn, msg.Content[0].(map[string]interface{}))
	case REGISTER:
		handleRegister(conn, msg.Content[0].(map[string]interface{}))
	case CHATMESSAGE:
		handleChatMessage(conn, msg.Content)
	case GROUPMESSAGE:
		handleGroupMessage(conn, msg.Content)
	}
}

func handleLogin(conn net.Conn, content map[string]interface{}) {
	log.Print("handleLogin...")
	msg := User{ID: content["id"].(string), Password: content["password"].(string)}
	log.Print(msg)

	var resMsg string
	var statusCode int

	if ok, err := msg.Login(); ok {
		// store string and IP in map
		userIP[msg.ID] = conn.RemoteAddr().String()
		log.Println("user ID: ", msg.ID, "\tIP: ", userIP[msg.ID])

		resMsg = "Login Success"
		statusCode = LOGIN

		// push messages to client if any
		// TODO

	} else if err != nil {
		resMsg = "User Not Found"
		statusCode = ERROR
		log.Print(err)
	} else {
		resMsg = "Password Not Match"
		statusCode = ERROR
		log.Print(err)
	}

	Response(conn, statusCode, resMsg)
}

func handleLogout(conn net.Conn, content map[string]interface{}) {
	log.Print("handleLogout...")
	msg := User{ID: content["id"].(string)}
	if _, ok := userIP[msg.ID]; ok {
		delete(userIP, msg.ID)

		log.Println("string: ", msg.ID, "\tIP: ", userIP[msg.ID])
		Response(conn, LOGOUT, "Logout Success")

		// close connection
		conn.Close()
	} else {
		Response(conn, ERROR, "User Not Found")
	}
}

func handleRegister(conn net.Conn, content map[string]interface{}) {
	// msg := User{ID: content["id"].(string), Name: content["name"].(string), Password: content["password"].(string)}
	msg := User{ID: content["id"].(string), Password: content["password"].(string)}
	if msg.CheckExist() {
		Response(conn, ERROR, "User Already Exist")
	} else {
		if err := msg.Register(); err != nil {
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
			// try to push message to client
			// TODO

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
			// try to push message to client
			//TODO

		} else {
			Response(conn, ERROR, "Sender Not Found")
		}
	}
}
