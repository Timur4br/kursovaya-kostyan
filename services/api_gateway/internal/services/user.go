package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type UserDAO interface {
	CreateUser(context.Context, map[string]any) (models.User, error)
	GetUser(context.Context, int64) (models.User, error)
	UpdateUser(context.Context, int64, map[string]any) (models.User, error)
	DeleteUser(context.Context, int64) error
}
type UserService struct{ dao UserDAO }

func NewUserService(dao UserDAO) *UserService { return &UserService{dao: dao} }

func (s *UserService) Create(ctx context.Context, input dto.UserInput) (models.User, error) {
	values, err := resources.Must("user").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.User{}, err
	}
	return s.dao.CreateUser(ctx, values)
}
func (s *UserService) Get(ctx context.Context, id int64) (models.User, error) {
	if err := validIDs(id); err != nil {
		return models.User{}, err
	}
	return s.dao.GetUser(ctx, id)
}
func (s *UserService) Update(ctx context.Context, id int64, input dto.UserInput) (models.User, error) {
	if err := validIDs(id); err != nil {
		return models.User{}, err
	}
	values, err := resources.Must("user").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.User{}, err
	}
	return s.dao.UpdateUser(ctx, id, values)
}
func (s *UserService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteUser(ctx, id)
}
