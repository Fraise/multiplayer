package message

import (
	"time"
)

type Message struct {
	Timestamp time.Time
	Players   []Player
}

type Player struct {
	Name     string
	Position Position
}

type Position struct {
	X float64
	Y float64
	Z float64
}
