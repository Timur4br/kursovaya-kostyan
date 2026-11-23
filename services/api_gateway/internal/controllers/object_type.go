package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type ObjectTypeService interface {
	Create(context.Context, dto.ObjectTypeInput) (models.ObjectType, error)
	Get(context.Context, int64) (models.ObjectType, error)
	Update(context.Context, int64, dto.ObjectTypeInput) (models.ObjectType, error)
	Delete(context.Context, int64) error
}
type ObjectTypeController struct{ service ObjectTypeService }

func NewObjectTypeController(service ObjectTypeService) *ObjectTypeController {
	return &ObjectTypeController{service: service}
}
func (c *ObjectTypeController) Create(ctx context.Context, input dto.ObjectTypeInput) (models.ObjectType, error) {
	return c.service.Create(ctx, input)
}
func (c *ObjectTypeController) Get(ctx context.Context, id int64) (models.ObjectType, error) {
	return c.service.Get(ctx, id)
}
func (c *ObjectTypeController) Update(ctx context.Context, id int64, input dto.ObjectTypeInput) (models.ObjectType, error) {
	return c.service.Update(ctx, id, input)
}
func (c *ObjectTypeController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
