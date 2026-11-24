package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type QueryHistoryService interface {
	Create(context.Context, dto.QueryHistoryInput) (models.QueryHistory, error)
	Get(context.Context, int64) (models.QueryHistory, error)
	Update(context.Context, int64, dto.QueryHistoryInput) (models.QueryHistory, error)
	Delete(context.Context, int64) error
}
type QueryHistoryController struct{ service QueryHistoryService }

func NewQueryHistoryController(service QueryHistoryService) *QueryHistoryController {
	return &QueryHistoryController{service: service}
}
func (c *QueryHistoryController) Create(ctx context.Context, input dto.QueryHistoryInput) (models.QueryHistory, error) {
	return c.service.Create(ctx, input)
}
func (c *QueryHistoryController) Get(ctx context.Context, id int64) (models.QueryHistory, error) {
	return c.service.Get(ctx, id)
}
func (c *QueryHistoryController) Update(ctx context.Context, id int64, input dto.QueryHistoryInput) (models.QueryHistory, error) {
	return c.service.Update(ctx, id, input)
}
func (c *QueryHistoryController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
