# anonymous-chatroom

## Feature

- Based on TCP protocal.
- Messages are encrypted by AES.
- Support emoji.

## Example

Server:

```bash
IP=0.0.0.0 PORT=9700 ./server
# 2021/06/06 21:04:02 listening on 0.0.0.0:9700
```

Client 1:

```bash
./netcat
# Please input server:
# 192.168.1.101:9700

# 本聊天室采用 AES 对称加密，无明文传输
# You are 123.112.10.143:18698
# no one is in this chatroom
```

Client 2:

```bash
./netcat
# Please input server:
# 192.168.1.101:9700

# 本聊天室采用 AES 对称加密，无明文传输
# You are 123.112.10.143:18852
# 123.112.10.143:18698 is in this chatroom too

# Client 1: 123.112.10.143:18852 has arrived
```

Client3:

```bash
./netcat
# Please input server:
# 192.168.1.101:9700

# 本聊天室采用 AES 对称加密，无明文传输
# You are 123.112.10.143:19395
# 123.112.10.143:18698 is in this chatroom too
# 123.112.10.143:18852 is in this chatroom too

# Client 1: 123.112.10.143:19395 has arrived
# Client 2: 123.112.10.143:19395 has arrived
```

## Usage

Build:

```bash
# server
go build server.go tcp.go cipher.go
# client
go build netcat.go tcp.go cipher.go
```

Server:

```bash
IP=0.0.0.0 PORT=9700 ./server
```

Client:

```bash
./netcat
# then input server address like 192.168.1.101:9700
```
