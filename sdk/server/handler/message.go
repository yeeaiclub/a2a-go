package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumosx/a2a-go/sdk/types"
)

type MessageRequestHandler struct {
	handler MessageHandler
}

func NewRequestHandler(handler MessageHandler) *MessageRequestHandler {
	return &MessageRequestHandler{handler: handler}
}

func (h *MessageRequestHandler) MessageRoute(engine *gin.Engine) {
	group := engine.Group("/message")
	group.POST("/send", h.OnMessageSend)
	group.POST("/stream", h.OnMessageSend)
}

func (h *MessageRequestHandler) OnMessageSend(ctx *gin.Context) {
	var message types.Message
	if err := ctx.Bind(&message); err != nil {
		ctx.JSON(http.StatusOK, types.JSONRPCErrorResponse[string]{
			Error: types.JSONParseError(err.Error()),
		})
	}

	message, err := h.handler.OnMessageSend(ctx, types.MessageSendParam{})
	if err != nil {
		ctx.JSON(http.StatusOK, "")
	}

	ctx.JSON(http.StatusOK, message)
}

func (h *MessageRequestHandler) OnMessageSendStream(ctx *gin.Context) {

}
