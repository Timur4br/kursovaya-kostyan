package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type UserService interface {
	Create(context.Context, dto.UserInput) (models.User, error)
	Get(context.Context, int64) (models.User, error)
	Update(context.Context, int64, dto.UserInput) (models.User, error)
	Delete(context.Context, int64) error
}
type UserController struct{ service UserService }

func NewUserController(service UserService) *UserController { return &UserController{service: service} }
func (c *UserController) Create(ctx context.Context, input dto.UserInput) (models.User, error) {
	return c.service.Create(ctx, input)
}
func (c *UserController) Get(ctx context.Context, id int64) (models.User, error) {
	return c.service.Get(ctx, id)
}
func (c *UserController) Update(ctx context.Context, id int64, input dto.UserInput) (models.User, error) {
	return c.service.Update(ctx, id, input)
}
func (c *UserController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
