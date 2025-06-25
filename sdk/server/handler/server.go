package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumosx/a2a-go/sdk/types"
)

type RequestHandler struct {
	basePath string
	port     int
}

func NewRequestHandler(basePath string) *RequestHandler {
	return &RequestHandler{basePath: basePath}
}

func (r *RequestHandler) Route(e *gin.Engine) {

}

func (r *RequestHandler) OnMessageSend(ctx *gin.Context) {
	var req types.SendMessageRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, types.JSONParseError(err.Error()))
		return
	}
}

func (r *RequestHandler) OnMessageSendStream(ctx *gin.Context) {
	var req types.SendMessageRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, types.JSONParseError(err.Error()))
		return
	}
}

func (r *RequestHandler) OnGetTask(ctx *gin.Context) {
	var req types.GetTaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, types.JSONParseError(err.Error()))
		return
	}
}

func (r *RequestHandler) OnCancelTask(ctx *gin.Context) {
	var req types.CancelTaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, types.JSONParseError(err.Error()))
		return
	}
}

func (r *RequestHandler) OnSetTaskPushNotificationConfig(ctx *gin.Context) {
	var req types.SetTaskPushNotificationConfigRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, types.JSONParseError(err.Error()))
		return
	}
}

func (r *RequestHandler) OnGetTaskPushNotificationConfig(ctx *gin.Context) {
	var req types.GetTaskPushNotificationConfigRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, types.JSONParseError(err.Error()))
		return
	}
}

func (r *RequestHandler) OnResubscribeToTask(ctx *gin.Context) {
	var req types.TaskResubscriptionRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, types.JSONParseError(err.Error()))
		return
	}
}
