package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type QueryEventController interface {
	Create(context.Context, dto.QueryEventInput) (models.QueryEvent, error)
	Get(context.Context, int64, int64) (models.QueryEvent, error)
	Update(context.Context, int64, int64, dto.QueryEventInput) (models.QueryEvent, error)
	Delete(context.Context, int64, int64) error
}
type QueryEventHandler struct {
	controller QueryEventController
	logger     *slog.Logger
}

func NewQueryEventHandler(controller QueryEventController, logger *slog.Logger) *QueryEventHandler {
	return &QueryEventHandler{controller: controller, logger: logger}
}

func (h *QueryEventHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.QueryEventInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/query_event/%d/%d", result.QueryID, result.EventID))
	c.JSON(201, gin.H{"data": result})
}
func (h *QueryEventHandler) Get(c *gin.Context) {
	queryID, ok := pathID(c, "query_id")
	if !ok {
		return
	}
	eventID, ok := pathID(c, "event_id")
	if !ok {
		return
	}
	result, err := h.controller.Get(c.Request.Context(), queryID, eventID)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.JSON(200, gin.H{"data": result})
}
func (h *QueryEventHandler) Update(c *gin.Context) {
	queryID, ok := pathID(c, "query_id")
	if !ok {
		return
	}
	eventID, ok := pathID(c, "event_id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.QueryEventInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), queryID, eventID, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/query_event/%d/%d", result.QueryID, result.EventID))
	c.JSON(200, gin.H{"data": result})
}
func (h *QueryEventHandler) Delete(c *gin.Context) {
	queryID, ok := pathID(c, "query_id")
	if !ok {
		return
	}
	eventID, ok := pathID(c, "event_id")
	if !ok {
		return
	}
	if err := h.controller.Delete(c.Request.Context(), queryID, eventID); err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Status(204)
}
