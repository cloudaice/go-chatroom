// Package server provides ...
package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func NewChatServer(bind string, rooms map[string]*Room, lock *sync.RWMutex) *ChatServer {
	return &ChatServer{bind, rooms, lock}
}

type ChatServer struct {
	Bind_to string
	Rooms   map[string]*Room
	lock    *sync.RWMutex
}

// GetRoom return a room, if this name of room is not exist,
// create a new room and return.
func (server *ChatServer) GetRoom(name string) *Room {
	server.lock.Lock()
	defer server.lock.Unlock()
	if _, ok := server.Rooms[name]; !ok {
		room := &Room{
			Server:  server,
			Name:    name,
			lock:    new(sync.RWMutex),
			Clients: make(map[string]*Client),
			In:      make(chan *Message),
		}
		go room.Listen()
		server.Rooms[name] = room
	}
	return server.Rooms[name]
}

// This method maybe should add a lock.
func (server *ChatServer) reportStatus() {
	for {
		time.Sleep(time.Second)
		server.lock.RLock()
		for _, room := range server.Rooms {
			fmt.Printf("Status: %s:%d\n", room.Name, len(room.Clients))
		}
		server.lock.RUnlock()
	}
}

func (server *ChatServer) ListenAndServe() {
	listener, err := net.Listen("tcp", server.Bind_to)

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	go server.reportStatus()
	// Main loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %s\n", err.Error())
			time.Sleep(time.Second)
			continue
		}

		c := &Client{
			Server: server,
			Name:   conn.RemoteAddr().String(),
			Conn:   conn,
			lock:   new(sync.RWMutex),
			Rooms:  make(map[string]*Room),
			In:     make(chan *Message, 100),
			Out:    make(chan *Message, 100),
			Quit:   make(chan struct{}),
		}
		go c.Listen()
		go c.Resp()
		go c.Recv()
	}
}
