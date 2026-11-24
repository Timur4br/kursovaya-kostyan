package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type LocationService interface {
	Create(context.Context, dto.LocationInput) (models.Location, error)
	Get(context.Context, int64) (models.Location, error)
	Update(context.Context, int64, dto.LocationInput) (models.Location, error)
	Delete(context.Context, int64) error
}
type LocationController struct{ service LocationService }

func NewLocationController(service LocationService) *LocationController {
	return &LocationController{service: service}
}
func (c *LocationController) Create(ctx context.Context, input dto.LocationInput) (models.Location, error) {
	return c.service.Create(ctx, input)
}
func (c *LocationController) Get(ctx context.Context, id int64) (models.Location, error) {
	return c.service.Get(ctx, id)
}
func (c *LocationController) Update(ctx context.Context, id int64, input dto.LocationInput) (models.Location, error) {
	return c.service.Update(ctx, id, input)
}
func (c *LocationController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
