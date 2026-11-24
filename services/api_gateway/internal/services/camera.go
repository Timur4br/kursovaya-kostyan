package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type CameraDAO interface {
	CreateCamera(context.Context, map[string]any) (models.Camera, error)
	GetCamera(context.Context, int64) (models.Camera, error)
	UpdateCamera(context.Context, int64, map[string]any) (models.Camera, error)
	DeleteCamera(context.Context, int64) error
}
type CameraService struct{ dao CameraDAO }

func NewCameraService(dao CameraDAO) *CameraService { return &CameraService{dao: dao} }

func (s *CameraService) Create(ctx context.Context, input dto.CameraInput) (models.Camera, error) {
	values, err := resources.Must("camera").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Camera{}, err
	}
	return s.dao.CreateCamera(ctx, values)
}
func (s *CameraService) Get(ctx context.Context, id int64) (models.Camera, error) {
	if err := validIDs(id); err != nil {
		return models.Camera{}, err
	}
	return s.dao.GetCamera(ctx, id)
}
func (s *CameraService) Update(ctx context.Context, id int64, input dto.CameraInput) (models.Camera, error) {
	if err := validIDs(id); err != nil {
		return models.Camera{}, err
	}
	values, err := resources.Must("camera").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Camera{}, err
	}
	return s.dao.UpdateCamera(ctx, id, values)
}
func (s *CameraService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteCamera(ctx, id)
}
