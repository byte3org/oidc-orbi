package storage

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

// SessionStore allows storing and retrieving sessions
type SessionStore struct {
	sessions map[string]webauthn.SessionData
}

// NewSessionStore initializes a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]webauthn.SessionData),
	}
}

// StoreSession stores a new session for a user
func (s *SessionStore) StoreSession(userID string, data webauthn.SessionData) {
	s.sessions[userID] = data
}

// GetSessionByUserID retrieves a session by a user ID
func (s *SessionStore) GetSessionByUserID(userID string) (webauthn.SessionData, bool) {
	data, ok := s.sessions[userID]
	return data, ok
}
