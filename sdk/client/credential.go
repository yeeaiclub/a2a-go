package auth

import (
	"errors"
	"sync"

	"github.com/yeeaiclub/a2a-go/sdk/client"
)

type Credential interface {
	GetCredentials(securitySchemeName string, context *client.CallContext) (string, error)
}

type CredentialKey struct {
	SessionID          string
	SecuritySchemeName string
}

type InMemoryContextCredentials struct {
	mu    sync.RWMutex
	store map[CredentialKey]string
}

func NewInMemoryContextCredentials() *InMemoryContextCredentials {
	return &InMemoryContextCredentials{
		store: make(map[CredentialKey]string),
	}
}

func (s *InMemoryContextCredentials) GetCredentials(securitySchemeName string, context *client.CallContext) (string, error) {
	if context == nil {
		return "", nil
	}
	sessionId, ok := context.State["sessionId"]
	if !ok {
		return "", nil
	}

	id, ok := sessionId.(string)
	if !ok {
		return "", errors.New("type asset failed")
	}

	key := CredentialKey{SessionID: id, SecuritySchemeName: securitySchemeName}
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func (s *InMemoryContextCredentials) SetCredentials(sessionId string, securitySchemeName, credential string) {
	key := CredentialKey{SessionID: sessionId, SecuritySchemeName: securitySchemeName}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = credential
}
