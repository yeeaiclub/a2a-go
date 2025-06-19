package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yumosx/a2a-go/sdk/types"
)

func TestGetAgentCard(t *testing.T) {
	testcases := []struct {
		name string
		want types.AgentCard
	}{
		{
			name: "success response",
			want: types.AgentCard{
				Name:               "hello, word",
				Description:        "a hello word agent",
				DefaultInputModes:  []string{"text"},
				DefaultOutputModes: []string{"text"},
				Version:            "1.0.0",
				Skills: []types.AgentSkill{
					{Id: "1", Description: "return hello", Name: "hello, word", Tags: []string{"hello, word"}},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				err := json.NewEncoder(w).Encode(&tc.want)
				require.NoError(t, err)
			}))
			resolver := NewA2ACardResolver(
				http.DefaultClient,
				server.URL,
			)
			card, err := resolver.GetAgentCard()
			require.NoError(t, err)
			assert.Equal(t, card, tc.want)
		})
	}
}
