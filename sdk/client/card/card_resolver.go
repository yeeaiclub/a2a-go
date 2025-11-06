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

package card

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/yeeaiclub/a2a-go/internal/logger"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// A2ACardResolver is responsible for retrieving agent card information from a specified endpoint.
// It provides a way to fetch [types.AgentCard] data using HTTP requests with configurable options.
type A2ACardResolver struct {
	client        *http.Client      // HTTP client used for making requests
	baseUrl       string            // Base URL of the agent service
	agentCardPath string            // Path to the agent card endpoint
	options       map[string]string // HTTP headers to include in requests
}

// NewA2ACardResolver creates a new instance of A2ACardResolver with the provided configuration.
func NewA2ACardResolver(client *http.Client, baseURL string, options ...A2ACardResolverOption) *A2ACardResolver {
	r := A2ACardResolver{
		client:        client,
		baseUrl:       baseURL,
		agentCardPath: types.AgentCardPath, // Default path from types package
	}

	// Apply all provided configuration options
	for _, opt := range options {
		r = opt.Option(r)
	}
	return &r
}

// GetAgentCard retrieves an agent card from the configured base URL.
// It sends an HTTP GET request to the agent card endpoint and returns the parsed card data.
func (a *A2ACardResolver) GetAgentCard(ctx context.Context) (*types.AgentCard, error) {
	targetURL, err := url.JoinPath(a.baseUrl, a.agentCardPath)
	if err != nil {
		return nil, fmt.Errorf("failed to construct url %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request %w", err)
	}

	if a.options != nil {
		for key, value := range a.options {
			req.Header.Set(key, value)
		}
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error(ctx, "failed to close response body", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get agent card: status code %d", resp.StatusCode)
	}

	var card types.AgentCard
	err = json.NewDecoder(resp.Body).Decode(&card)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}
	return &card, nil
}

// A2ACardResolverOption defines an interface for configuring A2ACardResolver instances.
type A2ACardResolverOption interface {
	Option(resolver A2ACardResolver) A2ACardResolver
}

type A2ACardResolverOptionFunc func(resolver A2ACardResolver) A2ACardResolver

// Option applies the function to the resolver and returns the updated resolver.
func (fn A2ACardResolverOptionFunc) Option(resolver A2ACardResolver) A2ACardResolver {
	return fn(resolver)
}

// WithAgentCardPath configures the specific path used to retrieve the agent card.
// By default, the resolver uses the standard path defined in types.AgentCardPath.
func WithAgentCardPath(path string) A2ACardResolverOption {
	return A2ACardResolverOptionFunc(func(resolver A2ACardResolver) A2ACardResolver {
		resolver.agentCardPath = path
		return resolver
	})
}

// WithHeader configures custom HTTP headers to be included in the agent card request.
func WithHeader(header map[string]string) A2ACardResolverOption {
	return A2ACardResolverOptionFunc(func(resolver A2ACardResolver) A2ACardResolver {
		resolver.options = header
		return resolver
	})
}
