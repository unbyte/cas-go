package cas

import (
	"sync"
	"time"
)

type Store interface {
	Get(sessionID string) (interface{}, bool)
	Set(sessionID string, data interface{}) error
	Del(sessionID string) error
}

type store struct {
	mu       sync.RWMutex
	store    map[string]interface{}
	oldStore map[string]interface{}
}

var _ Store = &store{}

func DefaultStore(clearInterval time.Duration) Store {
	if clearInterval < 0 {
		panic("clearInterval < 0 is invalid")
	}
	s := &store{
		store:    make(map[string]interface{}),
		oldStore: make(map[string]interface{}),
	}
	if clearInterval > 0 {
		go func() {
			timer := time.NewTicker(clearInterval)
			select {
			case <-timer.C:
				s.mu.Lock()
				s.oldStore = s.store
				s.store = make(map[string]interface{})
				s.mu.Unlock()
			}
		}()
	}
	return s
}

func (s *store) Get(sessionID string) (interface{}, bool) {
	s.mu.RLock()
	a, ok := s.store[sessionID]
	if !ok {
		a, ok = s.oldStore[sessionID]
	}
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
	delete(s.oldStore, sessionID)
	s.mu.Unlock()

	return nil
}
