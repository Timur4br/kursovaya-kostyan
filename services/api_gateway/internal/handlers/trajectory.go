package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type TrajectoryController interface {
	Create(context.Context, dto.TrajectoryInput) (models.Trajectory, error)
	Get(context.Context, int64) (models.Trajectory, error)
	Update(context.Context, int64, dto.TrajectoryInput) (models.Trajectory, error)
	Delete(context.Context, int64) error
}
type TrajectoryHandler struct {
	controller TrajectoryController
	logger     *slog.Logger
}

func NewTrajectoryHandler(controller TrajectoryController, logger *slog.Logger) *TrajectoryHandler {
	return &TrajectoryHandler{controller: controller, logger: logger}
}

func (h *TrajectoryHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.TrajectoryInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/trajectory/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *TrajectoryHandler) Get(c *gin.Context) {
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
func (h *TrajectoryHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.TrajectoryInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/trajectory/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *TrajectoryHandler) Delete(c *gin.Context) {
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
