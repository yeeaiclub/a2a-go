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

	log "github.com/yeeaiclub/a2a-go/internal/logger"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

const (
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 30 * time.Second
)

// Server implements the main HTTP server for agent APIs, including JSON-RPC and streaming endpoints.
type Server struct {
	agentCardPath string          // Path for agent card metadata
	card          types.AgentCard // Agent card metadata
	handler       Handler         // Business logic handler
	basePath      string          // Base path for API
	readTimeout   time.Duration   // HTTP read timeout
	writeTimeout  time.Duration   // HTTP write timeout
	idleTimeout   time.Duration   // HTTP idle timeout
}

// NewServer creates a new Server with the given configuration and options.
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

// Start launches the HTTP server on the specified port.
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
	log.Infof("Starting HTTP server on :%d with ReadTimeout=%v, WriteTimeout=%v, IdleTimeout=%v",
		port, s.readTimeout, s.writeTimeout, s.idleTimeout)
	return server.ListenAndServe()
}

// handleGetAgentCard handles GET requests for the agent card metadata.
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

// ServeHTTP is the main entry for JSON-RPC POST requests, dispatching to the appropriate handler.
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
		log.Warnf("Unknown method: %s", request.Method)
	}
}

// handleMessageSend handles the message/send JSON-RPC method.
func (s *Server) handleMessageSend(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	log.Infof("handleMessageSend called | id=%s, method=%s", id, request.Method)
	params, err := types.MapTo[types.MessageSendParam](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnMessageSend(ctx, params)
	if err != nil {
		log.Errorf("handleMessageSend | onMessageSend | %v", err)
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

// handleMessageSendStream handles the message/stream JSON-RPC method with server-sent events (SSE).
func (s *Server) handleMessageSendStream(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	log.Infof("handleMessageSendStream called | id=%s, method=%s", id, request.Method)
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
				return
			}
			switch ev.Type {
			case types.EventData:
				_ = ev.EncodeJSONRPC(encoder, id)
				flusher.Flush()
			case types.EventError:
				_ = ev.EncodeJSONRPC(encoder, id)
				flusher.Flush()
				return
			case types.EventDone:
				_ = ev.EncodeJSONRPC(encoder, id)
				flusher.Flush()
				return
			case types.EventClosed:
				return
			default:
			}
		}
	}
}

// handleGetTask handles the tasks/get JSON-RPC method.
func (s *Server) handleGetTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	log.Infof("handleGetTask called | id=%s, method=%s", id, request.Method)
	params, err := types.MapTo[types.TaskQueryParams](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnGetTask(ctx, params)
	if err != nil {
		log.Errorf("handleGetTask | onGetTask| %v", err)
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

// handleCancelTask handles the tasks/cancel JSON-RPC method.
func (s *Server) handleCancelTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	log.Infof("handleCancelTask called | id=%s, method=%s", id, request.Method)
	params, err := types.MapTo[types.TaskIdParams](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnCancelTask(ctx, params)
	if err != nil {
		log.Errorf("handleCancelTaskk | onCancelTask | %v", err)
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

// handleSetTaskPushNotificationConfig handles the tasks/pushNotificationConfig/set JSON-RPC method.
func (s *Server) handleSetTaskPushNotificationConfig(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	log.Infof("handleSetTaskPushNotificationConfig called | id=%s, method=%s", id, request.Method)
	params, err := types.MapTo[types.TaskPushNotificationConfig](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}
	event, err := s.handler.OnSetTaskPushNotificationConfig(ctx, params)
	if err != nil {
		log.Errorf("handleSetTaskPushNotificationConfig | OnSetTaskPushNotificationConfig | %v", err)
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

// handleGetTaskPushNotificationConfig handles the tasks/pushNotificationConfig/get JSON-RPC method.
func (s *Server) handleGetTaskPushNotificationConfig(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	log.Infof("handleGetTaskPushNotificationConfig called | id=%s, method=%s", id, request.Method)
	params, err := types.MapTo[types.TaskIdParams](request.Params)
	if err != nil {
		s.sendError(w, id, types.JSONParseError(err))
		return
	}

	event, err := s.handler.OnGetTaskPushNotificationConfig(ctx, params)
	if err != nil {
		log.Errorf("handleGetTaskPushNotificationConfig | OnGetTaskPushNotificationConfig | %v", err)
		s.sendError(w, id, types.InternalError())
		return
	}
	s.sendResponse(w, id, event)
}

// handleResubscribeToTask handles the tasks/resubscribe JSON-RPC method with SSE.
func (s *Server) handleResubscribeToTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
	log.Infof("handleResubscribeToTask called | id=%s, method=%s", id, request.Method)
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
				return
			}
			switch ev.Type {
			case types.EventData:
				_ = ev.EncodeJSONRPC(encoder, id)
				flusher.Flush()
			case types.EventError:
				_ = ev.EncodeJSONRPC(encoder, id)
				flusher.Flush()
				return
			case types.EventDone:
				_ = ev.EncodeJSONRPC(encoder, id)
				flusher.Flush()
				return
			case types.EventClosed:
				return
			default:
			}
		}
	}
}

// sendError writes a JSON-RPC error response.
func (s *Server) sendError(w http.ResponseWriter, id string, err *types.JSONRPCError) {
	response := types.JSONRPCResponse{
		Id:      id,
		JSONRPC: types.Version,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// sendResponse writes a JSON-RPC success response.
func (s *Server) sendResponse(w http.ResponseWriter, id string, result any) {
	response := types.JSONRPCResponse{
		Id:      id,
		JSONRPC: types.Version,
		Result:  result,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// ServerConfigOption allows customizing the Server via functional options.
type ServerConfigOption interface {
	Option(server *Server)
}

// ServerConfigOptionFunc is a function type for ServerConfigOption.
type ServerConfigOptionFunc func(server *Server)

func (fn ServerConfigOptionFunc) Option(server *Server) {
	fn(server)
}

// WithReadTimeout sets the server's read timeout.
func WithReadTimeout(readTimeout time.Duration) ServerConfigOption {
	return ServerConfigOptionFunc(func(server *Server) {
		server.readTimeout = readTimeout
	})
}

// WithWriteTimeout sets the server's write timeout.
func WithWriteTimeout(writeTimeout time.Duration) ServerConfigOption {
	return ServerConfigOptionFunc(func(server *Server) {
		server.writeTimeout = writeTimeout
	})
}

// WithIdleTimeout sets the server's idle timeout.
func WithIdleTimeout(idleTimeout time.Duration) ServerConfigOption {
	return ServerConfigOptionFunc(func(server *Server) {
		server.idleTimeout = idleTimeout
	})
}
