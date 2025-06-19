package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumosx/a2a-go/sdk/types"
)

type CardRequestHandler struct {
	card types.AgentCard
}

func NewCardRequestHandler(card types.AgentCard) *CardRequestHandler {
	return &CardRequestHandler{card: card}
}

func (c *CardRequestHandler) Route(engine *gin.Engine) {
	engine.GET("/get_card", c.GetCard)
}

func (c *CardRequestHandler) GetCard(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.card)
}
