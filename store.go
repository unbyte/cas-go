package cas

import (
	"sync"
)

type Store interface {
	Get(sessionID string) (interface{}, bool)
	Set(sessionID string, data interface{}) error
	Del(sessionID string) error
}

type store struct {
	mu    sync.RWMutex
	store map[string]interface{}
}

var _ Store = &store{}

func DefaultStore() Store {
	return &store{
		store: make(map[string]interface{}),
	}
}

func (s *store) Get(sessionID string) (interface{}, bool) {
	s.mu.RLock()
	a, ok := s.store[sessionID]
	s.mu.RUnlock()

	return a, ok
}

func (s *store) Set(sessionID string, attr interface{}) error {
	s.mu.Lock()
	s.store[sessionID] = attr
	s.mu.Unlock()

	return nil
}

func (s *store) Del(sessionID string) error {
	s.mu.Lock()
	delete(s.store, sessionID)
	s.mu.Unlock()

	return nil
}
