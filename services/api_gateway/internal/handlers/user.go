package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type UserController interface {
	Create(context.Context, dto.UserInput) (models.User, error)
	Get(context.Context, int64) (models.User, error)
	Update(context.Context, int64, dto.UserInput) (models.User, error)
	Delete(context.Context, int64) error
}
type UserHandler struct {
	controller UserController
	logger     *slog.Logger
}

func NewUserHandler(controller UserController, logger *slog.Logger) *UserHandler {
	return &UserHandler{controller: controller, logger: logger}
}

func (h *UserHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.UserInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/user/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *UserHandler) Get(c *gin.Context) {
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
func (h *UserHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.UserInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/user/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *UserHandler) Delete(c *gin.Context) {
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
