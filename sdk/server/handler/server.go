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

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yumosx/a2a-go/sdk/types"
)

type Server struct {
	card     types.AgentCard
	handler  Handler
	basePath string
}

func NewServer(card types.AgentCard, handler Handler, basePath string) *Server {
	return &Server{card: card, handler: handler, basePath: basePath}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()
	mux.Handle(s.basePath, s)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
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
		s.HandleMessageSend(r.Context(), w, &request, request.Id)
	case types.MethodMessageStream:
		s.HandleMessageSendStream(r.Context(), w, &request, request.Id)
	case types.MethodTasksGet:
		s.HandleGetTask(r.Context(), w, &request, request.Id)
	case types.MethodTasksCancel:
		s.HandleCancelTask(r.Context(), w, &request, request.Id)
	case types.MethodPushNotificationSet:
		s.HandleSetTaskPushNotificationConfig(r.Context(), w, &request, request.Id)
	case types.MethodPushNotificationGet:
		s.HandleGetTaskPushNotificationConfig(r.Context(), w, &request, request.Id)
	case types.MethodTasksResubscribe:
		s.HandleResubscribeToTask(r.Context(), w, &request, request.Id)
	}
}

func (s *Server) HandleMessageSend(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
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

func (s *Server) HandleMessageSendStream(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
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
		case ev, ok := <-events:
			if !ok {
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

func (s *Server) HandleGetTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
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

func (s *Server) HandleCancelTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
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

func (s *Server) HandleSetTaskPushNotificationConfig(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
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

func (s *Server) HandleGetTaskPushNotificationConfig(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
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

func (s *Server) HandleResubscribeToTask(ctx context.Context, w http.ResponseWriter, request *types.JSONRPCRequest, id string) {
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
		case ev, ok := <-events:
			if !ok {
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
