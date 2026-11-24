package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type TrajectoryPointController interface {
	Create(context.Context, dto.TrajectoryPointInput) (models.TrajectoryPoint, error)
	Get(context.Context, int64) (models.TrajectoryPoint, error)
	Update(context.Context, int64, dto.TrajectoryPointInput) (models.TrajectoryPoint, error)
	Delete(context.Context, int64) error
}
type TrajectoryPointHandler struct {
	controller TrajectoryPointController
	logger     *slog.Logger
}

func NewTrajectoryPointHandler(controller TrajectoryPointController, logger *slog.Logger) *TrajectoryPointHandler {
	return &TrajectoryPointHandler{controller: controller, logger: logger}
}

func (h *TrajectoryPointHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.TrajectoryPointInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/trajectory_point/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *TrajectoryPointHandler) Get(c *gin.Context) {
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
func (h *TrajectoryPointHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.TrajectoryPointInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/trajectory_point/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *TrajectoryPointHandler) Delete(c *gin.Context) {
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
