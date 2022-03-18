package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
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
