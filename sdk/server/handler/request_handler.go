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
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumosx/a2a-go/sdk/types"
)

type RequestHandler struct {
	handler Handler
}

func NewRequestHandler(handler Handler) *RequestHandler {
	return &RequestHandler{handler: handler}
}

func (r *RequestHandler) Route(e *echo.Echo) {
	e.POST("/", r.OnGetCard)
	e.POST("/message/send", r.OnMessageSend)
	e.POST("/message/stream", r.OnMessageSendStream)
	e.POST("/tasks/get", r.OnGetTask)
	e.POST("/tasks/cancel", r.OnCancelTask)
	e.POST("/tasks/resubscribe", r.OnResubscribeToTask)
	e.POST("/tasks/pushNotificationConfig/get", r.OnGetTaskPushNotificationConfig)
	e.POST("/tasks/pushNotificationConfig/set", r.OnSetTaskPushNotificationConfig)
}

func (r *RequestHandler) OnGetCard(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, r.handler.OnGetCard())
}

func (r *RequestHandler) OnMessageSend(ctx echo.Context) error {
	var req types.SendMessageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.JSONParseError(err)))
	}
	if req.Method != types.MethodMessageSend {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.MethodNotFoundError()))
	}
	event, err := r.handler.OnMessageSend(ctx.Request().Context(), req.Params)
	if err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.InternalError()))
	}
	return ctx.JSON(http.StatusOK, types.JSONRPCSuccessResponse(req.Id, event))
}

func (r *RequestHandler) OnMessageSendStream(ctx echo.Context) error {
	var req types.SendMessageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.JSONParseError(err)))
	}
	if req.Method != types.MethodMessageStream {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.MethodNotFoundError()))
	}

	w := ctx.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	events := r.handler.OnMessageSendStream(ctx.Request().Context(), req.Params)
	for {
		select {
		case <-ctx.Request().Context().Done():
			return ctx.Request().Context().Err()
		case event, ok := <-events:
			if !ok {
				return nil
			}
			err := event.MarshalTo(w, req.Id)
			if err != nil {
				return err
			}
			w.Flush()
			if event.Done() {
				return nil
			}
		}
	}
}

func (r *RequestHandler) OnGetTask(ctx echo.Context) error {
	var req types.GetTaskRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.JSONParseError(err)))
	}
	if req.Method != types.MethodTasksGet {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.MethodNotFoundError()))
	}
	task, err := r.handler.OnGetTask(ctx.Request().Context(), req.Params)
	if err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.InternalError()))
	}
	return ctx.JSON(http.StatusOK, types.JSONRPCSuccessResponse(req.Id, task))
}

func (r *RequestHandler) OnCancelTask(ctx echo.Context) error {
	var req types.CancelTaskRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.JSONParseError(err)))
	}
	if req.Method != types.MethodTasksCancel {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.MethodNotFoundError()))
	}

	task, err := r.handler.OnCancelTask(ctx.Request().Context(), req.Params)
	if err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.InternalError()))
	}
	return ctx.JSON(http.StatusOK, types.JSONRPCSuccessResponse(req.Id, task))
}

func (r *RequestHandler) OnSetTaskPushNotificationConfig(ctx echo.Context) error {
	var req types.SetTaskPushNotificationConfigRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.JSONParseError(err)))
	}
	if req.Method == types.MethodPushNotificationSet {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.MethodNotFoundError()))
	}
	config, err := r.handler.OnSetTaskPushNotificationConfig(ctx.Request().Context(), req.Params)
	if err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.InternalError()))
	}
	return ctx.JSON(http.StatusOK, types.JSONRPCSuccessResponse(req.Id, config))
}

func (r *RequestHandler) OnGetTaskPushNotificationConfig(ctx echo.Context) error {
	var req types.GetTaskPushNotificationConfigRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.JSONParseError(err)))
	}
	if req.Method == types.MethodPushNotificationGet {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.MethodNotFoundError()))
	}
	config, err := r.handler.OnGetTaskPushNotificationConfig(ctx.Request().Context(), req.Params)
	if err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.InternalError()))
	}
	return ctx.JSON(http.StatusOK, types.JSONRPCSuccessResponse(req.Id, config))
}

func (r *RequestHandler) OnResubscribeToTask(ctx echo.Context) error {
	var req types.TaskResubscriptionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.JSONParseError(err)))
	}
	if req.Method == types.MethodTasksResubscribe {
		return ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse(req.Id, types.MethodNotFoundError()))
	}

	events := r.handler.OnResubscribeToTask(ctx.Request().Context(), req.Params)
	w := ctx.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for {
		select {
		case <-ctx.Request().Context().Done():
			return ctx.Request().Context().Err()
		case event, ok := <-events:
			if !ok {
				return nil
			}
			err := event.MarshalTo(w, req.Id)
			if err != nil {
				return err
			}
			w.Flush()
			if event.Done() {
				return nil
			}
		}
	}
}
