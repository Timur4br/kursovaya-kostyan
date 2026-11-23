package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type EventService interface {
	Create(context.Context, dto.EventInput) (models.Event, error)
	Get(context.Context, int64) (models.Event, error)
	Update(context.Context, int64, dto.EventInput) (models.Event, error)
	Delete(context.Context, int64) error
}
type EventController struct{ service EventService }

func NewEventController(service EventService) *EventController {
	return &EventController{service: service}
}
func (c *EventController) Create(ctx context.Context, input dto.EventInput) (models.Event, error) {
	return c.service.Create(ctx, input)
}
func (c *EventController) Get(ctx context.Context, id int64) (models.Event, error) {
	return c.service.Get(ctx, id)
}
func (c *EventController) Update(ctx context.Context, id int64, input dto.EventInput) (models.Event, error) {
	return c.service.Update(ctx, id, input)
}
func (c *EventController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
