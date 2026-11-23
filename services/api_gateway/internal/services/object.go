package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type ObjectDAO interface {
	CreateObject(context.Context, map[string]any) (models.Object, error)
	GetObject(context.Context, int64) (models.Object, error)
	UpdateObject(context.Context, int64, map[string]any) (models.Object, error)
	DeleteObject(context.Context, int64) error
}
type ObjectService struct{ dao ObjectDAO }

func NewObjectService(dao ObjectDAO) *ObjectService { return &ObjectService{dao: dao} }

func (s *ObjectService) Create(ctx context.Context, input dto.ObjectInput) (models.Object, error) {
	values, err := resources.Must("object").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Object{}, err
	}
	return s.dao.CreateObject(ctx, values)
}
func (s *ObjectService) Get(ctx context.Context, id int64) (models.Object, error) {
	if err := validIDs(id); err != nil {
		return models.Object{}, err
	}
	return s.dao.GetObject(ctx, id)
}
func (s *ObjectService) Update(ctx context.Context, id int64, input dto.ObjectInput) (models.Object, error) {
	if err := validIDs(id); err != nil {
		return models.Object{}, err
	}
	values, err := resources.Must("object").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Object{}, err
	}
	return s.dao.UpdateObject(ctx, id, values)
}
func (s *ObjectService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteObject(ctx, id)
}
