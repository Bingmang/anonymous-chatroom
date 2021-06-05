package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

func packetSlitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 检查 atEOF 参数 和 数据包头部的四个字节是否 为 0x123456(我们定义的协议的魔数)
	if !atEOF && len(data) > 6 && binary.BigEndian.Uint32(data[:4]) == 0x123456 {
		var l int16
		// 读出 数据包中 实际数据 的长度(大小为 0 ~ 2^16)
		binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &l)
		pl := int(l) + 6
		if pl <= len(data) {
			return pl, data[:pl], nil
		}
	}
	return
}

func packet(data []byte) []byte {
	magicNum := make([]byte, 4)
	binary.BigEndian.PutUint32(magicNum, 0x123456)
	lenNum := make([]byte, 2)
	binary.BigEndian.PutUint16(lenNum, uint16(len(data)))
	packetBuf := bytes.NewBuffer(magicNum)
	packetBuf.Write(lenNum)
	packetBuf.Write(data)
	return packetBuf.Bytes()
}

func listenConn(conn net.Conn, tcpPacket chan []byte) {
	result := bytes.NewBuffer(nil)
	var buf [65542]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				close(tcpPacket)
				return
			} else {
				log.Println("tcp: listenConn: read error:", err)
			}
		}
		scanner := bufio.NewScanner(result)
		scanner.Split(packetSlitFunc)
		for scanner.Scan() {
			tcpPacket <- scanner.Bytes()[6:]
		}
		result.Reset()
	}
}
