package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type PersonService interface {
	Create(context.Context, dto.PersonInput) (models.Person, error)
	Get(context.Context, int64) (models.Person, error)
	Update(context.Context, int64, dto.PersonInput) (models.Person, error)
	Delete(context.Context, int64) error
}
type PersonController struct{ service PersonService }

func NewPersonController(service PersonService) *PersonController {
	return &PersonController{service: service}
}
func (c *PersonController) Create(ctx context.Context, input dto.PersonInput) (models.Person, error) {
	return c.service.Create(ctx, input)
}
func (c *PersonController) Get(ctx context.Context, id int64) (models.Person, error) {
	return c.service.Get(ctx, id)
}
func (c *PersonController) Update(ctx context.Context, id int64, input dto.PersonInput) (models.Person, error) {
	return c.service.Update(ctx, id, input)
}
func (c *PersonController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
