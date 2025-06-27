// Copyright 2025 yumosx
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

package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yumosx/a2a-go/sdk/types"
)

type A2ACardResolver struct {
	client        *http.Client
	baseUrl       string
	agentCardPath string
	options       map[string]string
}

type A2ACardResolverOption interface {
	Option(resolver A2ACardResolver) A2ACardResolver
}

type A2ACardResolverOptionFunc func(resolver A2ACardResolver) A2ACardResolver

func (fn A2ACardResolverOptionFunc) Option(resolver A2ACardResolver) A2ACardResolver {
	return fn(resolver)
}

func WithAgentCardPath(agentCardPath string) A2ACardResolverOption {
	return A2ACardResolverOptionFunc(func(resolver A2ACardResolver) A2ACardResolver {
		resolver.agentCardPath = agentCardPath
		return resolver
	})
}
func WithOptions(options map[string]string) A2ACardResolverOption {
	return A2ACardResolverOptionFunc(func(resolver A2ACardResolver) A2ACardResolver {
		resolver.options = options
		return resolver
	})
}

func NewA2ACardResolver(client *http.Client, baseUrl string, options ...A2ACardResolverOption) *A2ACardResolver {
	r := A2ACardResolver{
		client:        client,
		baseUrl:       baseUrl,
		agentCardPath: types.AgentCardPath,
	}

	for _, opt := range options {
		r = opt.Option(r)
	}
	return &r
}

func (a *A2ACardResolver) GetAgentCard() (types.AgentCard, error) {
	targetUrl := fmt.Sprintf("%s/%s", a.baseUrl, a.agentCardPath)
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return types.AgentCard{}, err
	}
	if a.options != nil {
		for key, value := range a.options {
			req.Header.Set(key, value)
		}
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return types.AgentCard{}, err
	}
	var card types.AgentCard
	err = json.NewDecoder(resp.Body).Decode(&card)
	if err != nil {
		return types.AgentCard{}, err
	}
	return card, nil
}
