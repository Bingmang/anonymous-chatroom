package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string)
	users = make(map[string]bool)
)

func main() {
	ip := os.Getenv("IP")
	if ip == "" {
		ip = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9700"
	}
	listenAddress := ip + ":" + port
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on", listenAddress)
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	ch := make(chan string)
	go clientWriter(conn, ch)

	// 欢迎消息
	who := conn.RemoteAddr().String()
	ch <- "本聊天室采用 AES 对称加密，无明文传输"
	ch <- "You are " + who
	if len(users) == 0 {
		ch <- "no one is in this chatroom"
	} else {
		for addr := range users {
			ch <- addr + " is in this chatroom too"
		}
	}
	users[who] = true

	// 入场通知
	messages <- who + " has arrived"
	entering <- ch

	// 检测输入
	tcpPacket := make(chan []byte)
	go listenConn(conn, tcpPacket)
	for bytes := range tcpPacket {
		decryptedMsg, err := decrypt(bytes)
		if err != nil {
			log.Println(err)
		}
		messages <- who + ": " + string(decryptedMsg)
	}

	// 离场通知
	leaving <- ch
	messages <- who + " has left"
	delete(users, who)
}

// 消息出口，所有出口消息都加密了
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Println("DEBUG: clientWriter: msg", msg)
		encryptedMsg, err := encrypt([]byte(msg))
		if err != nil {
			fmt.Println(err)
		}
		conn.Write(packet(encryptedMsg))
	}
}
