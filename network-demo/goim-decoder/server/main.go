package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"

	"goim-decoder/pkg"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
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
	reader := bufio.NewReader(c)
	for {
		peek, err := reader.Peek(pkg.PackageLengthSize())
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				break
			} else {
				log.Println("ending...")
				log.Println()
			}
			break
		}
		buffer := bytes.NewBuffer(peek)
		var size int32
		if err := binary.Read(buffer, binary.BigEndian, &size); err != nil {
			log.Println(err)
		}
		if int32(reader.Buffered()) < size {
			continue
		}
		data := make([]byte, size)
		if _, err := reader.Read(data); err != nil {
			log.Println(err)
			continue
		}
		content, err := pkg.Decoder(data)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(string(content.Content))
	}
}
