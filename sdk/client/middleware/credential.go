// Copyright 2025 yeeaiclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yeeaiclub/a2a-go/sdk/web"
)

type Credential interface {
	GetCredentials(securitySchemeName string, context web.Context) (string, error)
}

type CredentialKey struct {
	SessionID          string
	SecuritySchemeName string
}

type InMemoryContextCredentials struct {
	mu sync.RWMutex
	//  "sessionId" -> id and schemeName -> token
	store map[CredentialKey]string
}

func NewInMemoryContextCredentials() *InMemoryContextCredentials {
	return &InMemoryContextCredentials{
		store: make(map[CredentialKey]string),
	}
}

func (s *InMemoryContextCredentials) GetCredentials(securitySchemeName string, context web.Context) (string, error) {
	if context == nil {
		return "", nil
	}
	sessionId := context.Get("sessionId")
	if sessionId == nil {
		return "", nil
	}
	id, ok := sessionId.(string)
	if !ok {
		return "", fmt.Errorf("sessionId type assert failed: got %T, want string", sessionId)
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

func (s *InMemoryContextCredentials) Get(scheme string, context web.Context) (string, error) {
	if context == nil {
		return "", nil
	}
	sessionId := context.Get("sessionId")

	id, ok := sessionId.(string)
	if !ok {
		return "", errors.New("type asset failed")
	}

	key := CredentialKey{SessionID: id, SecuritySchemeName: scheme}
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
