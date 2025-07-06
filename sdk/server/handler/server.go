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

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yeeaiclub/a2a-go/sdk/types"
)

const (
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 30 * time.Second
)

type Server struct {
	agentCardPath string
	card          types.AgentCard
	handler       Handler
	basePath      string
	readTimeout   time.Duration
	writeTimeout  time.Duration
	idleTimeout   time.Duration
}

func NewServer(cardPath string, basePath string, card types.AgentCard, handler Handler, options ...ServerConfigOption) *Server {
	server := &Server{
		basePath:      basePath,
		agentCardPath: cardPath,
		card:          card,
		handler:       handler,
		readTimeout:   defaultReadTimeout,
		writeTimeout:  defaultWriteTimeout,
		idleTimeout:   defaultIdleTimeout,
	}
	for _, opt := range options {
		opt.Option(server)
	}
	return server
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc(s.agentCardPath, s.handleGetAgentCard)
	mux.Handle(s.basePath, s)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}
	return server.ListenAndServe()
}

func (s *Server) handleGetAgentCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(s.card); err != nil {
		s.sendError(w, "", types.JSONParseError(err))
		return
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request types.JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.sendError(w, request.Id, types.JSONParseError(err))
		return
	}
	switch request.Method {
	case types.MethodMessageSend:
		s.handleMessageSend(r.Context(), w, &request, request.Id)
	case types.MethodMessageStream:
		s.handleMessageSendStream(r.Context(), w, &request, request.Id)
	case types.MethodTasksGet:
		s.handleGetTask(r.Context(), w, &request, request.Id)
	case types.MethodTasksCancel:
		s.handleCancelTask(r.Context(), w, &request, request.Id)
	case types.MethodPushNotificationSet:
		s.handleSetTaskPushNotificationConfig(r.Context(), w, &request, request.Id)
	case types.MethodPushNotificationGet:
		s.handleGetTaskPushNotificationConfig(r.Context(), w, &request, request.Id)
	case types.MethodTasksResubscribe:
		s.handleResubscribeToTask(r.Context(), w, &request, request.Id)
	default:
	}
}

func (s *Server) handleMessageSend(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	params, err := types.MapTo[types.MessageSendParam](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnMessageSend(ctx, params)
	if err != nil {
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

func (s *Server) handleMessageSendStream(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	params, err := types.MapTo[types.MessageSendParam](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	events := s.handler.OnMessageSendStream(ctx, params)

	for {
		select {
		case <-ctx.Done():
			return
		case ev, o := <-events:
			if !o {
				if ev.Err != nil || ev.Event != nil {
					_ = ev.EncodeJSONRPC(encoder, id)
				}
				return
			}
			if ev.Err != nil {
				_ = ev.EncodeJSONRPC(encoder, id)
				return
			}
			err = ev.EncodeJSONRPC(encoder, id)
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func (s *Server) handleGetTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	params, err := types.MapTo[types.TaskQueryParams](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnGetTask(ctx, params)
	if err != nil {
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

func (s *Server) handleCancelTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	params, err := types.MapTo[types.TaskIdParams](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnCancelTask(ctx, params)
	if err != nil {
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

func (s *Server) handleSetTaskPushNotificationConfig(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	params, err := types.MapTo[types.TaskPushNotificationConfig](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnSetTaskPushNotificationConfig(ctx, params)
	if err != nil {
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

func (s *Server) handleGetTaskPushNotificationConfig(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	params, err := types.MapTo[types.TaskIdParams](request.Params)
	if err != nil {
		s.sendError(w, id, types.InternalError())
		return
	}

	event, err := s.handler.OnGetTaskPushNotificationConfig(ctx, params)
	if err != nil {
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

func (s *Server) handleResubscribeToTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	params, err := types.MapTo[types.TaskIdParams](request.Params)
	if err != nil {
		s.sendError(w, id, types.InternalError())
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	events := s.handler.OnResubscribeToTask(ctx, params)

	for {
		select {
		case <-ctx.Done():
			return
		case ev, o := <-events:
			if !o {
				if ev.Err != nil || ev.Event != nil {
					_ = ev.EncodeJSONRPC(encoder, id)
				}
				return
			}
			if ev.Err != nil {
				_ = ev.EncodeJSONRPC(encoder, id)
				return
			}
			err = ev.EncodeJSONRPC(encoder, id)
			if err != nil {
				return
			}
			flusher.Flush()
			if ev.Done() {
				return
			}
		}
	}
}

func (s *Server) sendError(w http.ResponseWriter, id string, err *types.JSONRPCError) {
	response := types.JSONRPCResponse{
		Id:      id,
		JSONRPC: types.Version,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) sendResponse(w http.ResponseWriter, id string, result any) {
	response := types.JSONRPCResponse{
		Id:      id,
		JSONRPC: types.Version,
		Result:  result,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

type ServerConfigOption interface {
	Option(server *Server)
}

type ServerConfigOptionFunc func(server *Server)

func (fn ServerConfigOptionFunc) Option(server *Server) {
	fn(server)
}

func WithReadTimeout(readTimeout time.Duration) ServerConfigOption {
	return ServerConfigOptionFunc(func(server *Server) {
		server.readTimeout = readTimeout
	})
}

func WithWriteTimeout(writeTimeout time.Duration) ServerConfigOption {
	return ServerConfigOptionFunc(func(server *Server) {
		server.writeTimeout = writeTimeout
	})
}

func WithIdleTimeout(idleTimeout time.Duration) ServerConfigOption {
	return ServerConfigOptionFunc(func(server *Server) {
		server.idleTimeout = idleTimeout
	})
}
