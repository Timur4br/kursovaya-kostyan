package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type CameraController interface {
	Create(context.Context, dto.CameraInput) (models.Camera, error)
	Get(context.Context, int64) (models.Camera, error)
	Update(context.Context, int64, dto.CameraInput) (models.Camera, error)
	Delete(context.Context, int64) error
}
type CameraHandler struct {
	controller CameraController
	logger     *slog.Logger
}

func NewCameraHandler(controller CameraController, logger *slog.Logger) *CameraHandler {
	return &CameraHandler{controller: controller, logger: logger}
}

func (h *CameraHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.CameraInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/camera/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *CameraHandler) Get(c *gin.Context) {
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
func (h *CameraHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.CameraInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/camera/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *CameraHandler) Delete(c *gin.Context) {
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
