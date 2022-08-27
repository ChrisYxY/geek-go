package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

const bufferSize = 1024

func main() {
	after := time.After(5 * time.Second)
	iter := 0
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// sending message
	for {
		select {
		case <-after:
			log.Println("time out")
			return
		default:
			for i := 0; i < 10; i++ {
				content := "hello[" + strconv.Itoa(iter) + "]"
				_, err = conn.Write(patch(content))
				if err != nil {
					log.Fatal(err)
				}
				iter++
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// patching empty byte into origin message
func patch(message string) []byte {
	res := make([]byte, bufferSize)
	copy(res, []byte(message))
	return res
}
