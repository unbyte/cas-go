package cas

import (
	"github.com/unbyte/cas-go/parser"
	"sync"
)

type SessionStore interface {
	Get(sessionID string) (parser.Attributes, bool)
	Set(sessionID string, attributes parser.Attributes) error
	Del(sessionID string) error
}

type sessionStore struct {
	mu    sync.RWMutex
	store map[string]parser.Attributes
}

var _ SessionStore = &sessionStore{}

func DefaultSessionStore() SessionStore {
	return &sessionStore{
		store: make(map[string]parser.Attributes),
	}
}

func (s *sessionStore) Get(sessionID string) (parser.Attributes, bool) {
	s.mu.RLock()
	a, ok := s.store[sessionID]
	s.mu.RUnlock()

	return a, ok
}

func (s *sessionStore) Set(sessionID string, attr parser.Attributes) error {
	s.mu.Lock()
	s.store[sessionID] = attr
	s.mu.Unlock()
	return nil
}

func (s *sessionStore) Del(sessionID string) error {
	s.mu.Lock()
	delete(s.store, sessionID)
	s.mu.Unlock()

	return nil
}
