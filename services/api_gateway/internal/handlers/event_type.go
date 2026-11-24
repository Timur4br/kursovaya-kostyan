package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type EventTypeController interface {
	Create(context.Context, dto.EventTypeInput) (models.EventType, error)
	Get(context.Context, int64) (models.EventType, error)
	Update(context.Context, int64, dto.EventTypeInput) (models.EventType, error)
	Delete(context.Context, int64) error
}
type EventTypeHandler struct {
	controller EventTypeController
	logger     *slog.Logger
}

func NewEventTypeHandler(controller EventTypeController, logger *slog.Logger) *EventTypeHandler {
	return &EventTypeHandler{controller: controller, logger: logger}
}

func (h *EventTypeHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.EventTypeInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/event_type/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *EventTypeHandler) Get(c *gin.Context) {
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
func (h *EventTypeHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.EventTypeInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/event_type/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *EventTypeHandler) Delete(c *gin.Context) {
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
