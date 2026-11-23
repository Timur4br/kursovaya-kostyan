package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type VideoSegmentController interface {
	Create(context.Context, dto.VideoSegmentInput) (models.VideoSegment, error)
	Get(context.Context, int64) (models.VideoSegment, error)
	Update(context.Context, int64, dto.VideoSegmentInput) (models.VideoSegment, error)
	Delete(context.Context, int64) error
}
type VideoSegmentHandler struct {
	controller VideoSegmentController
	logger     *slog.Logger
}

func NewVideoSegmentHandler(controller VideoSegmentController, logger *slog.Logger) *VideoSegmentHandler {
	return &VideoSegmentHandler{controller: controller, logger: logger}
}

func (h *VideoSegmentHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.VideoSegmentInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/video_segment/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *VideoSegmentHandler) Get(c *gin.Context) {
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
func (h *VideoSegmentHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.VideoSegmentInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/video_segment/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *VideoSegmentHandler) Delete(c *gin.Context) {
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
