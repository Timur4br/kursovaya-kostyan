package handlers

import (
	"context"
	"fmt"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type PersonController interface {
	Create(context.Context, dto.PersonInput) (models.Person, error)
	Get(context.Context, int64) (models.Person, error)
	Update(context.Context, int64, dto.PersonInput) (models.Person, error)
	Delete(context.Context, int64) error
}
type PersonHandler struct {
	controller PersonController
	logger     *slog.Logger
}

func NewPersonHandler(controller PersonController, logger *slog.Logger) *PersonHandler {
	return &PersonHandler{controller: controller, logger: logger}
}

func (h *PersonHandler) Create(c *gin.Context) {
	input, ok := bindJSON[dto.PersonInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/person/%d", result.ID))
	c.JSON(201, gin.H{"data": result})
}
func (h *PersonHandler) Get(c *gin.Context) {
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
func (h *PersonHandler) Update(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	input, ok := bindJSON[dto.PersonInput](c)
	if !ok {
		return
	}
	result, err := h.controller.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err, h.logger)
		return
	}
	c.Header("Location", fmt.Sprintf("/person/%d", result.ID))
	c.JSON(200, gin.H{"data": result})
}
func (h *PersonHandler) Delete(c *gin.Context) {
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
