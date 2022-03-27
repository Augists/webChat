package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
)

func GetRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetMD5(str string) string {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func RemoveSlice(slice []string, i string) []string {
	for j, v := range slice {
		if v == i {
			slice = append(slice[:j], slice[j+1:]...)
			break
		}
	}
	return slice
}

func Response(conn net.Conn, statusCode int, content string) {
	resContent := []interface{}{content}
	res := clientMessageAPI{Type: statusCode, Content: resContent}
	resJSON, _ := json.Marshal(res)
	conn.Write(resJSON)
	fmt.Println("Response:", string(resJSON))
}

func ListToString(list []string) string {
	var str string
	for _, v := range list {
		str += v + ","
	}
	return str[:len(str)-1]
}

func StringToList(str string) []string {
	return strings.Split(str, ",")
}
