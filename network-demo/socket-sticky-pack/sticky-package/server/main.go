package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

const bufferSize = 1024

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listing from 0.0.0.0:8080")

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handler(conn)
	}
}

func handler(c net.Conn) {
	defer c.Close()
	buf := make([]byte, bufferSize)
	result := bytes.NewBuffer(nil)
	for {
		n, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("ending: " + err.Error())
				return
			} else {
				log.Println("Read error: " + err.Error())
				break
			}
		}
		result.Write(buf[0:n])
		log.Printf("recevie size[%d]: %v", n, result.String())
		result.Reset()
	}
}
