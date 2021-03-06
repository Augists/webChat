package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"sync"
)

var (
	userIP map[string]string
	wg     sync.WaitGroup
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
	/*
	 * use bytes.Buffer to store message instead of []byte
	 */
	buf := bytes.Buffer{}
	// buf := bytes.NewBuffer(nil)
	// var buf [1024]byte
	/*
	 * expand buffer size to 1024
	 */
	buf.Grow(1024)
	for {
		n, err := buf.ReadFrom(conn)
		// n, err := conn.Read(buf.Bytes())
		// n, err := conn.Read(buf[0:])
		if err != nil {
			log.Print(err)
		} else {
			if n == 0 {
				continue
			} else {
				/*
				 * Jsonify in python will use <Response>
				 * but json.dumps will use <string> with " surrounded
				 * so we won't deliver " in buf to handleMessage
				 */
				go handleMessage(conn, buf.Bytes()[1:n-1])
				// go handleMessage(conn, buf[1:n-1])
				// go handleMessage(conn, buf[0:n])
			}
		}
		/*
		 * flush buffer to empty after each message
		 */
		buf.Reset()
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
	u := User{ID: content["id"].(string), Password: content["password"].(string)}
	log.Print(u)

	if ok, err := u.Login(); ok {
		/*
		 * store string and IP in map
		 */
		userIP[u.ID] = conn.RemoteAddr().String()
		log.Println("user ID: ", u.ID, "\tIP: ", userIP[u.ID])

		resMsg := "Login Success"
		log.Print(resMsg)
		Response(conn, LOGIN, resMsg)

		/*
		 * push messages to client if any
		 */
		wg.Add(2)
		go u.PullMessage(conn)
		go u.PullGroupMessage(conn)
		/*
		 * waitGroup until pull complete
		 */
		wg.Wait()
	} else {
		resMsg := err.Error()
		log.Print(resMsg)
		Response(conn, ERROR, resMsg)
	}
}

func handleLogout(conn net.Conn, content map[string]interface{}) {
	log.Print("handleLogout...")
	msg := User{ID: content["id"].(string)}
	if _, ok := userIP[msg.ID]; ok {
		delete(userIP, msg.ID)

		log.Println("string: ", msg.ID, "\tIP: ", userIP[msg.ID])
		resMsg := "Logout Success"
		Response(conn, LOGOUT, resMsg)
		log.Print("Logout Success")

		/*
		 * close connection and goroutine
		 */
		conn.Close()
		return
	} else {
		resMsg := "User Not Found"
		Response(conn, ERROR, resMsg)
		log.Print(resMsg)
	}
}

func handleRegister(conn net.Conn, content map[string]interface{}) {
	log.Print("handleRegister...")
	msg := User{ID: content["id"].(string), Name: content["name"].(string), Password: content["password"].(string)}
	// msg := User{ID: content["id"].(string), Password: content["password"].(string)}
	if msg.CheckExist() {
		resMsg := "User Already Exist"
		log.Print(resMsg)
		Response(conn, ERROR, resMsg)
	} else {
		if err := msg.Register(); err != nil {
			resMsg := "Register Failed"
			log.Print(resMsg)
			log.Print(err)
			Response(conn, ERROR, resMsg)
		} else {
			resMsg := "Register Success"
			log.Print(resMsg)
			Response(conn, REGISTER, resMsg)
		}
	}
}

func handleChatMessage(conn net.Conn, content []interface{}) {
	msg := []Message{}
	var onlineCount int = 0
	var offlineCount int = 0
	for i, v := range content {
		json.Unmarshal(v.([]byte), &msg[i])
		if _, ok := userIP[msg[i].ReceiverID]; ok {
			onlineCount++
			/*
			 * receiver online
			 * try to push message to client
			 */
			clientConn, err := net.Dial("tcp", userIP[msg[i].SenderID])
			if err != nil {
				log.Print(err)
			} else {
				clientConn.Write(v.([]byte))
				clientConn.Close()
			}
		} else {
			/*
			 * receiver offline
			 * store in database
			 */
			offlineCount++
			msg[i].Store()
		}
	}
	Response(conn, CHATMESSAGE, "Online User:"+strconv.Itoa(onlineCount)+",Offline User:"+strconv.Itoa(offlineCount))
}

func handleGroupMessage(conn net.Conn, content []interface{}) {
	msg := []GroupMessage{}
	for i, v := range content {
		json.Unmarshal(v.([]byte), &msg[i])
		/*
		 * try to push message to client
		 */
		gid := msg[i].GroupID
		groupUserList := GetGroupUserList(gid)
		for _, userID := range groupUserList {
			if _, ok := userIP[userID]; ok {
				/*
				 * receiver online
				 * try to push message to client
				 */
				clientConn, err := net.Dial("tcp", userIP[userID])
				if err != nil {
					log.Print(err)
				} else {
					clientConn.Write(v.([]byte))
					clientConn.Close()
				}
			} else {
				/*
				 * receiver offline
				 * store in database
				 */
				msg[i].Store()
			}
		}
	}
	Response(conn, GROUPMESSAGE, "")
}
