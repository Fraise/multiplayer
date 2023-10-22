package state

import (
	"errors"
	"fmt"
	"multiplayer/message"
	"multiplayer/player"
	"sync"
)

type State struct {
	connMut  *sync.RWMutex
	sessions map[string]*sessionState

	maxSessions uint
	maxPlayers  uint

	currentSessionsNum uint
}

func New(options Options) *State {
	return &State{
		connMut:            new(sync.RWMutex),
		sessions:           make(map[string]*sessionState, options.MaxSessions),
		maxSessions:        options.MaxSessions,
		maxPlayers:         options.MaxPlayer,
		currentSessionsNum: 0,
	}
}

func (s *State) AddPlayer(sessionId string, player *player.Player) error {
	s.connMut.Lock()
	defer s.connMut.Unlock()

	// if the session already exists
	if session, ok := s.sessions[sessionId]; ok {
		if session.currentPlayersNum >= s.maxPlayers {
			return errors.New("too many players in that session")
		}

		session.addPlayer(player)

		return nil
	}

	// if the session doesn't exist, check the max number of sessions first
	if s.currentSessionsNum >= s.maxSessions {
		return errors.New("too many sessions")
	}

	// create a new session and add the player
	session := &sessionState{}
	session.addPlayer(player)
	s.sessions[sessionId] = session

	return nil
}

func (s *State) GetSessionStateMessage(sessionId string) (message.Message, error) {
	if session, ok := s.sessions[sessionId]; ok {
		return session.ToMessage(), nil
	}

	return message.Message{}, errors.New("invalid session id")
}

func (s *State) RemovePlayer(sessionId string, p *player.Player) error {
	s.connMut.Lock()
	defer s.connMut.Unlock()

	if session, ok := s.sessions[sessionId]; ok {
		_, err := session.removePlayer(p)

		if err != nil {
			return fmt.Errorf("could not remove player from state: %w", err)
		}

		return nil
	}

	return errors.New("session does not exist")
}
