package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumosx/a2a-go/sdk/types"
)

type TaskRequestHandler struct {
	handler TaskHandler
}

func NewTaskHandler(handler TaskHandler) TaskRequestHandler {
	return TaskRequestHandler{handler: handler}
}

func (task *TaskRequestHandler) Route(engine *gin.Engine) {
	group := engine.Group("/tasks")
	group.POST("/get", task.OnGetTask)
	group.POST("/cancel", task.OnCancelTask)
	group.POST("/pushNotificationConfig/set", task.OnSetTaskPushNotificationConfig)
	group.POST("/pushNotificationConfig/get", task.OnGetTaskPushNotificationConfig)
	group.POST("/tasks/resubscribe", task.OnResubscribeToTask)
}

func (task *TaskRequestHandler) OnGetTask(ctx *gin.Context) {
	var req types.GetTaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse[string]{
			Error: types.JSONParseError(err.Error()),
		})
		return
	}
	t, err := task.handler.OnGetTask(ctx, types.TaskQueryParam{})
	if err != nil {
		ctx.JSON(http.StatusOK, "")
	}
	ctx.JSON(http.StatusOK, t)
}

func (task *TaskRequestHandler) OnCancelTask(ctx *gin.Context) {
	var req types.CancelTaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse[string]{
			Error: types.JSONParseError(err.Error()),
		})
		return
	}
	t, err := task.handler.OnCancelTask(ctx, types.TaskIdParams{})
	if err != nil {
		ctx.JSON(http.StatusOK, "")
	}
	ctx.JSON(http.StatusOK, t)
}

func (task *TaskRequestHandler) OnSetTaskPushNotificationConfig(ctx *gin.Context) {
	var req types.SetTaskPushNotificationConfigRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse[string]{
			Error: types.JSONParseError(err.Error()),
		})
		return
	}
	t, err := task.handler.OnSetTaskPushNotificationConfig(ctx, types.TaskPushNotificationConfig{})
	if err != nil {
		ctx.JSON(http.StatusOK, "")
		return
	}
	ctx.JSON(http.StatusOK, t)
}

func (task *TaskRequestHandler) OnGetTaskPushNotificationConfig(ctx *gin.Context) {
	var req types.GetTaskPushNotificationConfigRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse[string]{
			Error: types.JSONParseError(err.Error()),
		})
		return
	}
	t, err := task.handler.OnGetTaskPushNotificationConfig(ctx, types.TaskIdParams{})
	if err != nil {
		ctx.JSON(http.StatusOK, "")
		return
	}
	ctx.JSON(http.StatusOK, t)
}

func (task *TaskRequestHandler) OnResubscribeToTask(ctx *gin.Context) {
	var req types.TaskResubscriptionRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse[string]{
			Error: types.JSONParseError(err.Error()),
		})
		return
	}
}
