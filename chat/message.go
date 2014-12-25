// Package message provides ...
package chat

import (
	"time"
)

const (
	_ = iota
	NORMAL
	QUIT
	JOIN
	DISMISS
	PAUSE
	KICK
)

type Message struct {
	Sender *Client
	// room name
	Receiver string
	// command
	Command int
	Content interface{}
	Time    time.Time
}
