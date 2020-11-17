package cas

import (
	"net/http"
)

type SessionManager interface {
	Set(w http.ResponseWriter, r *http.Request, sessionID string) error
	Get(r *http.Request) (string, bool)
}

type sessionManager struct {
	cookieName string
}

var _ SessionManager = &sessionManager{}

func (s *sessionManager) Set(w http.ResponseWriter, r *http.Request, sessionID string) error {
	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
	})
	return nil
}

func (s *sessionManager) Get(r *http.Request) (string, bool) {
	session, err := r.Cookie(s.cookieName)
	if err != nil {
		return "", false
	}
	return session.Value, true

}

func DefaultSessionManager(cookieName string) SessionManager {
	if cookieName == "" {
		cookieName = "cas-go"
	}
	return &sessionManager{cookieName: cookieName}
}
