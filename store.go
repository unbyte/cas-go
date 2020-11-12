package cas

import (
	"github.com/unbyte/cas-go/parser"
	"sync"
)

type SessionStore interface {
	Get(session string) (string, bool)
	Set(session, ticket string) error
	Del(session string) error
}

type TicketStore interface {
	Get(ticket string) (parser.Attributes, bool)
	Set(ticket string, attributes parser.Attributes) error
	Del(ticket string) error
}

type ticketStore struct {
	mu    sync.RWMutex
	store map[string]parser.Attributes
}

var _ TicketStore = &ticketStore{}

func DefaultTicketStore() TicketStore {
	return &ticketStore{
		mu:    sync.RWMutex{},
		store: map[string]parser.Attributes{},
	}
}

func (s *ticketStore) Get(ticket string) (parser.Attributes, bool) {
	s.mu.RLock()

	t, ok := s.store[ticket]
	s.mu.RUnlock()

	if !ok {
		return nil, false
	}

	return t, true
}

func (s *ticketStore) Set(ticket string, attr parser.Attributes) error {
	s.mu.Lock()

	s.store[ticket] = attr

	s.mu.Unlock()
	return nil
}

func (s *ticketStore) Del(ticket string) error {
	s.mu.Lock()
	delete(s.store, ticket)
	s.mu.Unlock()
	return nil
}

type sessionStore struct {
	mu    sync.RWMutex
	store map[string]string
}

var _ SessionStore = &sessionStore{}

func DefaultSessionStore() SessionStore {
	return &sessionStore{
		mu:    sync.RWMutex{},
		store: map[string]string{},
	}
}

func (s *sessionStore) Get(session string) (string, bool) {
	s.mu.RLock()
	t, ok := s.store[session]
	s.mu.RUnlock()

	if !ok {
		return "", false
	}

	return t, true
}

func (s *sessionStore) Set(session, ticket string) error {
	s.mu.Lock()
	s.store[session] = ticket
	s.mu.Unlock()
	return nil
}

func (s *sessionStore) Del(session string) error {
	s.mu.Lock()
	delete(s.store, session)
	s.mu.Unlock()
	return nil
}
