package cas

import (
	"github.com/unbyte/cas-go/parser"
	"sync"
)

type AttributesStore interface {
	Get(sessionID string) (parser.Attributes, bool)
	Set(sessionID string, attributes parser.Attributes) error
	Del(sessionID string) error
}

type attributesStore struct {
	mu    sync.RWMutex
	store map[string]parser.Attributes
}

var _ AttributesStore = &attributesStore{}

func DefaultSessionStore() AttributesStore {
	return &attributesStore{
		store: make(map[string]parser.Attributes),
	}
}

func (s *attributesStore) Get(sessionID string) (parser.Attributes, bool) {
	s.mu.RLock()
	a, ok := s.store[sessionID]
	s.mu.RUnlock()

	return a, ok
}

func (s *attributesStore) Set(sessionID string, attr parser.Attributes) error {
	s.mu.Lock()
	s.store[sessionID] = attr
	s.mu.Unlock()
	return nil
}

func (s *attributesStore) Del(sessionID string) error {
	s.mu.Lock()
	delete(s.store, sessionID)
	s.mu.Unlock()

	return nil
}
