package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func RandomString(length int) string {
	var result string
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz "
	for i := 0; i < length; i++ {
		result += string(chars[rand.Intn(len(chars))])
	}
	return result + "\n"
}

func main() {
	f, err := os.OpenFile("/wamedia/gofile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < 1029200; i++ {
		if i%1000 == 0 {
			fmt.Println("records generated : " + strconv.Itoa(i))
		}
		f.Write([]byte(RandomString(100)))
	}
}

