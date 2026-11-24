package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type EventController interface {
	Create(context.Context, dto.EventInput) (models.Event, error)
	Get(context.Context, int64) (models.Event, error)
	Update(context.Context, int64, dto.EventInput) (models.Event, error)
	Delete(context.Context, int64) error
}
type EventHandler struct {
	controller EventController
	logger     *slog.Logger
}

func NewEventHandler(controller EventController, logger *slog.Logger) *EventHandler {
	return &EventHandler{controller: controller, logger: logger}
}

func (h *EventHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.EventInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/event/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *EventHandler) Get(c *gin.Context) {
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
func (h *EventHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.EventInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/event/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *EventHandler) Delete(c *gin.Context) {
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
