package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type QueryHistoryController interface {
	Create(context.Context, dto.QueryHistoryInput) (models.QueryHistory, error)
	Get(context.Context, int64) (models.QueryHistory, error)
	Update(context.Context, int64, dto.QueryHistoryInput) (models.QueryHistory, error)
	Delete(context.Context, int64) error
}
type QueryHistoryHandler struct {
	controller QueryHistoryController
	logger     *slog.Logger
}

func NewQueryHistoryHandler(controller QueryHistoryController, logger *slog.Logger) *QueryHistoryHandler {
	return &QueryHistoryHandler{controller: controller, logger: logger}
}

func (h *QueryHistoryHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.QueryHistoryInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/query_history/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *QueryHistoryHandler) Get(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	result, err := h.controller.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.JSON(200, gin.H{"data": result})
}
func (h *QueryHistoryHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.QueryHistoryInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/query_history/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *QueryHistoryHandler) Delete(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	if err := h.controller.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Status(204)
}
