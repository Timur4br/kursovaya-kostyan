package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type LocationDAO interface {
	CreateLocation(context.Context, map[string]any) (models.Location, error)
	GetLocation(context.Context, int64) (models.Location, error)
	UpdateLocation(context.Context, int64, map[string]any) (models.Location, error)
	DeleteLocation(context.Context, int64) error
}
type LocationService struct{ dao LocationDAO }

func NewLocationService(dao LocationDAO) *LocationService { return &LocationService{dao: dao} }

func (s *LocationService) Create(ctx context.Context, input dto.LocationInput) (models.Location, error) {
	values, err := resources.Must("location").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Location{}, err
	}
	return s.dao.CreateLocation(ctx, values)
}
func (s *LocationService) Get(ctx context.Context, id int64) (models.Location, error) {
	if err := validIDs(id); err != nil {
		return models.Location{}, err
	}
	return s.dao.GetLocation(ctx, id)
}
func (s *LocationService) Update(ctx context.Context, id int64, input dto.LocationInput) (models.Location, error) {
	if err := validIDs(id); err != nil {
		return models.Location{}, err
	}
	values, err := resources.Must("location").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Location{}, err
	}
	return s.dao.UpdateLocation(ctx, id, values)
}
func (s *LocationService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteLocation(ctx, id)
}
