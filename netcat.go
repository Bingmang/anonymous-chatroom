package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)


func main() {
	fmt.Println("Please input server:")
	var address string
	fmt.Scan(&address)
	if address == "" {
		address = "localhost:9700"
	}
	conn, err := net.Dial("tcp", address)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}
	tcpPacket := make(chan []byte)
	go listenConn(conn, tcpPacket)
	go sendMessage(conn, bufio.NewScanner(os.Stdin))
	for bytes := range tcpPacket {
		decryptedBytes, err := decrypt(bytes)
		if err != nil {
			log.Println(err)
		}
		os.Stdout.Write(append(decryptedBytes, '\n'))
	}
	log.Println("服务器关闭")
}

// 消息入口，所有入口消息都加密了
func sendMessage(conn io.Writer, scanner *bufio.Scanner) {
	for scanner.Scan() {
		var message string
		message = scanner.Text()
		if len(message) == 0 {
			continue
		}
		encryptedMsg, err := encrypt([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write(packet(encryptedMsg))
		if err != nil {
			log.Fatal(err)
		}
	}
}
