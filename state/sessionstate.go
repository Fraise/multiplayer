package state

import (
	"errors"
	"multiplayer/message"
	"multiplayer/player"
	"time"
)

type sessionState struct {
	players           map[*player.Player]bool
	currentPlayersNum uint
}

func (s *sessionState) addPlayer(p *player.Player) uint {
	if s.players == nil {
		s.players = make(map[*player.Player]bool)
	}

	s.players[p] = true
	s.currentPlayersNum++

	return s.currentPlayersNum
}

func (s *sessionState) removePlayer(p *player.Player) (uint, error) {
	if s.players == nil {
		return 0, errors.New("no player in this session")
	}

	if _, ok := s.players[p]; !ok {
		return 0, errors.New("this player is not in a session")
	}

	delete(s.players, p)
	s.currentPlayersNum--

	return s.currentPlayersNum, nil
}

func (s *sessionState) ToMessage() message.Message {
	msg := message.Message{
		Timestamp: time.Now(),
		Players:   make([]message.Player, 0),
	}

	for p, _ := range s.players {
		msg.Players = append(msg.Players, p.ToMessage())
	}

	return msg
}
