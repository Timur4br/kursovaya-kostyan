package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type ObjectTypeDAO interface {
	CreateObjectType(context.Context, map[string]any) (models.ObjectType, error)
	GetObjectType(context.Context, int64) (models.ObjectType, error)
	UpdateObjectType(context.Context, int64, map[string]any) (models.ObjectType, error)
	DeleteObjectType(context.Context, int64) error
}
type ObjectTypeService struct{ dao ObjectTypeDAO }

func NewObjectTypeService(dao ObjectTypeDAO) *ObjectTypeService { return &ObjectTypeService{dao: dao} }

func (s *ObjectTypeService) Create(ctx context.Context, input dto.ObjectTypeInput) (models.ObjectType, error) {
	values, err := resources.Must("object_type").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.ObjectType{}, err
	}
	return s.dao.CreateObjectType(ctx, values)
}
func (s *ObjectTypeService) Get(ctx context.Context, id int64) (models.ObjectType, error) {
	if err := validIDs(id); err != nil {
		return models.ObjectType{}, err
	}
	return s.dao.GetObjectType(ctx, id)
}
func (s *ObjectTypeService) Update(ctx context.Context, id int64, input dto.ObjectTypeInput) (models.ObjectType, error) {
	if err := validIDs(id); err != nil {
		return models.ObjectType{}, err
	}
	values, err := resources.Must("object_type").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.ObjectType{}, err
	}
	return s.dao.UpdateObjectType(ctx, id, values)
}
func (s *ObjectTypeService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteObjectType(ctx, id)
}
