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

package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	log "github.com/yeeaiclub/a2a-go/internal/logger"
	"github.com/yeeaiclub/a2a-go/sdk/types"
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

func (c *A2AClient) SendMessage(params types.MessageSendParam) (*types.JSONRPCResponse, error) {
	req := types.SendMessageRequest{
		Id:     uuid.New().String(),
		Method: types.MethodMessageSend,
		Params: params,
	}
	var resp types.JSONRPCResponse
	err := c.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *A2AClient) GetTask(params types.TaskQueryParams) (*types.JSONRPCResponse, error) {
	req := types.GetTaskRequest{
		Id:     uuid.New().String(),
		Method: types.MethodTasksGet,
		Params: params,
	}

	var resp types.JSONRPCResponse
	err := c.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *A2AClient) CancelTask(params types.TaskIdParams) (*types.JSONRPCResponse, error) {
	req := types.CancelTaskRequest{
		Id:     uuid.New().String(),
		Method: types.MethodTasksCancel,
		Params: params,
	}
	var resp types.JSONRPCResponse
	err := c.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *A2AClient) SetTaskPushNotificationConfig(params types.TaskPushNotificationConfig) (*types.JSONRPCResponse, error) {
	req := types.SetTaskPushNotificationConfigRequest{
		Id:     uuid.New().String(),
		Method: types.MethodPushNotificationSet,
		Params: params,
	}

	var resp types.JSONRPCResponse
	err := c.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *A2AClient) GetTaskPushNotificationConfig(params types.TaskIdParams) (*types.JSONRPCResponse, error) {
	req := types.GetTaskPushNotificationConfigRequest{
		Id:     uuid.New().String(),
		Method: types.MethodPushNotificationGet,
		Params: params,
	}

	var resp types.JSONRPCResponse
	err := c.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *A2AClient) SendMessageStream(param types.MessageSendParam, eventChan chan any) error {
	request := types.SendStreamingMessageRequest{
		Id:      uuid.New().String(),
		JSONRPC: types.Version,
		Method:  types.MethodMessageStream,
		Params:  param,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	httpResp, err := c.clint.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		err = httpResp.Body.Close()
		if err != nil {
			log.Errorf("Failed to send HTTP request to %s: %v", c.url, err)
		}
	}()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}
	return c.processStream(httpReq.Context(), httpResp.Body, eventChan)
}

func (c *A2AClient) ResubscribeToTask(params types.TaskIdParams, eventChan chan any) error {
	request := types.TaskResubscriptionRequest{
		Id:      uuid.New().String(),
		JSONRPC: types.Version,
		Method:  types.MethodTasksResubscribe,
		Params:  params,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	httpResp, err := c.clint.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		err = httpResp.Body.Close()
		if err != nil {
			log.Errorf("Failed to send HTTP request to %s: %v", c.url, err)
		}
	}()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}
	return c.processStream(httpReq.Context(), httpResp.Body, eventChan)
}

func (c *A2AClient) sendRequest(request any, resp *types.JSONRPCResponse) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	httpResp, err := c.clint.Do(httpReq)
	if err != nil {
		return err
	}

	defer func() {
		err = httpResp.Body.Close()
		if err != nil {
			log.Errorf("Failed to send HTTP request to %s: %v", c.url, err)
		}
	}()

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return err
	}
	return nil
}

func (c *A2AClient) processStream(ctx context.Context, body io.Reader, eventChan chan any) error {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var event types.JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return fmt.Errorf("failed to decode event: %v", err)
		}
		if event.Error != nil {
			return fmt.Errorf("a2a error: %s (code: %d)", event.Error.Message, event.Error.Code)
		}
		result, err := json.Marshal(event.Result)
		if err != nil {
			return fmt.Errorf("failed to encode event result: %w", err)
		}
		select {
		case eventChan <- json.RawMessage(result):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	return nil
}
