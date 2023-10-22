package player

import (
	"multiplayer/message"
)

type Player struct {
	Name string
	Position
}

type Position struct {
	X float64
	Y float64
	Z float64
}

func (p *Player) ToMessage() message.Player {
	return message.Player{
		Name: p.Name,
		Position: message.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
		},
	}
}
