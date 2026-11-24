package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type LocationController interface {
	Create(context.Context, dto.LocationInput) (models.Location, error)
	Get(context.Context, int64) (models.Location, error)
	Update(context.Context, int64, dto.LocationInput) (models.Location, error)
	Delete(context.Context, int64) error
}
type LocationHandler struct {
	controller LocationController
	logger     *slog.Logger
}

func NewLocationHandler(controller LocationController, logger *slog.Logger) *LocationHandler {
	return &LocationHandler{controller: controller, logger: logger}
}

func (h *LocationHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.LocationInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/location/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *LocationHandler) Get(c *gin.Context) {
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
func (h *LocationHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.LocationInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/location/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *LocationHandler) Delete(c *gin.Context) {
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
