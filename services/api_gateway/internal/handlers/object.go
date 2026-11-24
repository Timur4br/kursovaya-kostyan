package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type ObjectController interface {
	Create(context.Context, dto.ObjectInput) (models.Object, error)
	Get(context.Context, int64) (models.Object, error)
	Update(context.Context, int64, dto.ObjectInput) (models.Object, error)
	Delete(context.Context, int64) error
}
type ObjectHandler struct {
	controller ObjectController
	logger     *slog.Logger
}

func NewObjectHandler(controller ObjectController, logger *slog.Logger) *ObjectHandler {
	return &ObjectHandler{controller: controller, logger: logger}
}

func (h *ObjectHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.ObjectInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/object/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *ObjectHandler) Get(c *gin.Context) {
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
func (h *ObjectHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.ObjectInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/object/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *ObjectHandler) Delete(c *gin.Context) {
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
