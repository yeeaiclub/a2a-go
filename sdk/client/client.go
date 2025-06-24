package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/yumosx/a2a-go/sdk/types"
)

type A2AClient struct {
	card  *types.AgentCard
	clint *http.Client
	url   string
}

type A2AClientOption interface {
	Option(client A2AClient) A2AClient
}

type A2AClientOptionFunc func(client A2AClient) A2AClient

func (fn A2AClientOptionFunc) Option(client A2AClient) A2AClient {
	return fn(client)
}

func WithUrl(url string) A2AClientOption {
	return A2AClientOptionFunc(func(client A2AClient) A2AClient {
		client.url = url
		return client
	})
}

func WithAgentCard(card *types.AgentCard) A2AClientOption {
	return A2AClientOptionFunc(func(client A2AClient) A2AClient {
		client.card = card
		return client
	})
}

func NewClient(client *http.Client, options ...A2AClientOption) (*A2AClient, error) {
	a2aClient := A2AClient{
		clint: client,
	}

	for _, opt := range options {
		a2aClient = opt.Option(a2aClient)
	}

	if a2aClient.url == "" && a2aClient.card == nil {
		return nil, errors.New("must provide either agent_card or url")
	}

	if a2aClient.card != nil {
		a2aClient.url = a2aClient.card.URL
	}
	return &a2aClient, nil
}

func (c *A2AClient) SendMessage(
	ctx context.Context,
	request types.SendMessageRequest,
	options map[string]string,
) (types.SendMessageResponse, error) {
	if request.Id == "" {
		request.Id = uuid.New().String()
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return types.SendMessageResponse{}, err
	}

	resp, err := c.sendRequest(ctx, payload, options)
	if err != nil {
		return types.SendMessageResponse{}, err
	}

	var response types.SendMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return types.SendMessageResponse{}, err
	}

	return response, nil
}

func (c *A2AClient) SendMessageStream(ctx context.Context, request types.SendStreamingMessageRequest) error {
	return nil
}

func (c *A2AClient) GetTask(
	ctx context.Context,
	request types.GetTaskRequest,
	options map[string]string,
) (types.GetTaskSuccessResponse, error) {
	if request.Id == "" {
		request.Id = uuid.New().String()
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return types.GetTaskSuccessResponse{}, err
	}
	resp, err := c.sendRequest(ctx, payload, options)
	if err != nil {
		return types.GetTaskSuccessResponse{}, err
	}
	var response types.GetTaskSuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return types.GetTaskSuccessResponse{}, err
	}
	return response, nil
}

func (c *A2AClient) CancelTask(
	ctx context.Context,
	request types.CancelTaskRequest,
	options map[string]string,
) (types.CancelTaskResponse, error) {
	if request.Id == "" {
		request.Id = uuid.New().String()
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return types.CancelTaskResponse{}, err
	}
	resp, err := c.sendRequest(ctx, payload, options)
	if err != nil {
		return types.CancelTaskResponse{}, err
	}
	var response types.CancelTaskResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return types.CancelTaskResponse{}, err
	}
	return response, nil
}

func (c *A2AClient) sendRequest(ctx context.Context, payload []byte, options map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	for key, value := range options {
		req.Header.Set(key, value)
	}

	resp, err := c.clint.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
