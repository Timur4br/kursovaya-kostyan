package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type ObjectTypeController interface {
	Create(context.Context, dto.ObjectTypeInput) (models.ObjectType, error)
	Get(context.Context, int64) (models.ObjectType, error)
	Update(context.Context, int64, dto.ObjectTypeInput) (models.ObjectType, error)
	Delete(context.Context, int64) error
}
type ObjectTypeHandler struct {
	controller ObjectTypeController
	logger     *slog.Logger
}

func NewObjectTypeHandler(controller ObjectTypeController, logger *slog.Logger) *ObjectTypeHandler {
	return &ObjectTypeHandler{controller: controller, logger: logger}
}

func (h *ObjectTypeHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.ObjectTypeInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/object_type/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *ObjectTypeHandler) Get(c *gin.Context) {
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
func (h *ObjectTypeHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.ObjectTypeInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/object_type/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *ObjectTypeHandler) Delete(c *gin.Context) {
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
