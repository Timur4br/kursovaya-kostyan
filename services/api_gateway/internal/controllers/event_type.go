package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type EventTypeService interface {
	Create(context.Context, dto.EventTypeInput) (models.EventType, error)
	Get(context.Context, int64) (models.EventType, error)
	Update(context.Context, int64, dto.EventTypeInput) (models.EventType, error)
	Delete(context.Context, int64) error
}
type EventTypeController struct{ service EventTypeService }

func NewEventTypeController(service EventTypeService) *EventTypeController {
	return &EventTypeController{service: service}
}
func (c *EventTypeController) Create(ctx context.Context, input dto.EventTypeInput) (models.EventType, error) {
	return c.service.Create(ctx, input)
}
func (c *EventTypeController) Get(ctx context.Context, id int64) (models.EventType, error) {
	return c.service.Get(ctx, id)
}
func (c *EventTypeController) Update(ctx context.Context, id int64, input dto.EventTypeInput) (models.EventType, error) {
	return c.service.Update(ctx, id, input)
}
func (c *EventTypeController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
