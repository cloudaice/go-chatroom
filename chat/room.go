// Package main provides ...
package chat

import (
	"fmt"
	"sync"
)

func kickName(msg *Message) string {
	return fmt.Sprintf("%s", msg.Content)
}

type Room struct {
	Server  *ChatServer
	Name    string
	lock    *sync.RWMutex
	Clients map[string]*Client
	In      chan *Message
}

func (r *Room) Listen() {
	fmt.Printf("Chatroom: %s opened\n", r.Name)
	for msg := range r.In {
		switch msg.Command {
		case QUIT:
			r.lock.Lock()
			delete(r.Clients, msg.Sender.Name)
			r.lock.Unlock()
			r.broadcast(msg)
		case JOIN:
			fmt.Printf("%s joined\n", msg.Sender.Name)
			r.lock.Lock()
			r.Clients[msg.Sender.Name] = msg.Sender
			r.lock.Unlock()

			r.broadcast(msg)
		case KICK:
			name := kickName(msg)
			r.lock.RLock()
			_, ok := r.Clients[name]
			r.lock.RUnlock()
			if ok {
				r.lock.Lock()
				delete(r.Clients, name)
				r.lock.Unlock()

				r.broadcast(msg)
			}
		case DISMISS:
			// Blocking broadcasting...
			// do nothing.
			r.broadcast(msg)
			return
		default:
			r.broadcast(msg)
		}
	}
}

func (r *Room) broadcast(msg *Message) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, c := range r.Clients {
		c.In <- msg
	}
}
