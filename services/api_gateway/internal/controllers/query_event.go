package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type QueryEventService interface {
	Create(context.Context, dto.QueryEventInput) (models.QueryEvent, error)
	Get(context.Context, int64, int64) (models.QueryEvent, error)
	Update(context.Context, int64, int64, dto.QueryEventInput) (models.QueryEvent, error)
	Delete(context.Context, int64, int64) error
}
type QueryEventController struct{ service QueryEventService }

func NewQueryEventController(service QueryEventService) *QueryEventController {
	return &QueryEventController{service: service}
}
func (c *QueryEventController) Create(ctx context.Context, input dto.QueryEventInput) (models.QueryEvent, error) {
	return c.service.Create(ctx, input)
}
func (c *QueryEventController) Get(ctx context.Context, queryID int64, eventID int64) (models.QueryEvent, error) {
	return c.service.Get(ctx, queryID, eventID)
}
func (c *QueryEventController) Update(ctx context.Context, queryID int64, eventID int64, input dto.QueryEventInput) (models.QueryEvent, error) {
	return c.service.Update(ctx, queryID, eventID, input)
}
func (c *QueryEventController) Delete(ctx context.Context, queryID int64, eventID int64) error {
	return c.service.Delete(ctx, queryID, eventID)
}
