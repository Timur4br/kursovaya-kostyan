package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type TrajectoryPointDAO interface {
	CreateTrajectoryPoint(context.Context, map[string]any) (models.TrajectoryPoint, error)
	GetTrajectoryPoint(context.Context, int64) (models.TrajectoryPoint, error)
	UpdateTrajectoryPoint(context.Context, int64, map[string]any) (models.TrajectoryPoint, error)
	DeleteTrajectoryPoint(context.Context, int64) error
}
type TrajectoryPointService struct{ dao TrajectoryPointDAO }

func NewTrajectoryPointService(dao TrajectoryPointDAO) *TrajectoryPointService {
	return &TrajectoryPointService{dao: dao}
}

func (s *TrajectoryPointService) Create(ctx context.Context, input dto.TrajectoryPointInput) (models.TrajectoryPoint, error) {
	values, err := resources.Must("trajectory_point").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.TrajectoryPoint{}, err
	}
	return s.dao.CreateTrajectoryPoint(ctx, values)
}
func (s *TrajectoryPointService) Get(ctx context.Context, id int64) (models.TrajectoryPoint, error) {
	if err := validIDs(id); err != nil {
		return models.TrajectoryPoint{}, err
	}
	return s.dao.GetTrajectoryPoint(ctx, id)
}
func (s *TrajectoryPointService) Update(ctx context.Context, id int64, input dto.TrajectoryPointInput) (models.TrajectoryPoint, error) {
	if err := validIDs(id); err != nil {
		return models.TrajectoryPoint{}, err
	}
	values, err := resources.Must("trajectory_point").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.TrajectoryPoint{}, err
	}
	return s.dao.UpdateTrajectoryPoint(ctx, id, values)
}
func (s *TrajectoryPointService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteTrajectoryPoint(ctx, id)
}
