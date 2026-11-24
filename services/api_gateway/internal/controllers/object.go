package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type ObjectService interface {
	Create(context.Context, dto.ObjectInput) (models.Object, error)
	Get(context.Context, int64) (models.Object, error)
	Update(context.Context, int64, dto.ObjectInput) (models.Object, error)
	Delete(context.Context, int64) error
}
type ObjectController struct{ service ObjectService }

func NewObjectController(service ObjectService) *ObjectController {
	return &ObjectController{service: service}
}
func (c *ObjectController) Create(ctx context.Context, input dto.ObjectInput) (models.Object, error) {
	return c.service.Create(ctx, input)
}
func (c *ObjectController) Get(ctx context.Context, id int64) (models.Object, error) {
	return c.service.Get(ctx, id)
}
func (c *ObjectController) Update(ctx context.Context, id int64, input dto.ObjectInput) (models.Object, error) {
	return c.service.Update(ctx, id, input)
}
func (c *ObjectController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
