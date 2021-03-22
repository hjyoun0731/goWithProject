package main

import (
	"github.com/stretchr/objx"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hjyoun0731/goWithProject/trace"
)

type room struct {
	// forward는 수신 메시지를 보관하는 채널이며
	// 수신한 메시지는 다른 클라이언트로 전달돼야 한다
	forward chan *message

	// join은 방에 들어오려는 클라이언트를 위한 채널이다
	join chan *client

	// leave는 방을 나가길 원하는 클라이언트를 위한 채널이다
	leave chan *client

	// clients는 현재 채팅방에 있는 모든 클라이언트를 보유한다
	clients map[*client]bool

	// tracer는 방 안에서 활동의 추적 정보를 수신한다.
	tracer trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 입장
			r.clients[client] = true
		case client := <-r.leave:
			// 퇴장
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			// 모든 클라이언트에게 메시지 전달
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
