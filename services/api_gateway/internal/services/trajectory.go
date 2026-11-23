package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type TrajectoryDAO interface {
	CreateTrajectory(context.Context, map[string]any) (models.Trajectory, error)
	GetTrajectory(context.Context, int64) (models.Trajectory, error)
	UpdateTrajectory(context.Context, int64, map[string]any) (models.Trajectory, error)
	DeleteTrajectory(context.Context, int64) error
}
type TrajectoryService struct{ dao TrajectoryDAO }

func NewTrajectoryService(dao TrajectoryDAO) *TrajectoryService { return &TrajectoryService{dao: dao} }

func (s *TrajectoryService) Create(ctx context.Context, input dto.TrajectoryInput) (models.Trajectory, error) {
	values, err := resources.Must("trajectory").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Trajectory{}, err
	}
	return s.dao.CreateTrajectory(ctx, values)
}
func (s *TrajectoryService) Get(ctx context.Context, id int64) (models.Trajectory, error) {
	if err := validIDs(id); err != nil {
		return models.Trajectory{}, err
	}
	return s.dao.GetTrajectory(ctx, id)
}
func (s *TrajectoryService) Update(ctx context.Context, id int64, input dto.TrajectoryInput) (models.Trajectory, error) {
	if err := validIDs(id); err != nil {
		return models.Trajectory{}, err
	}
	values, err := resources.Must("trajectory").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Trajectory{}, err
	}
	return s.dao.UpdateTrajectory(ctx, id, values)
}
func (s *TrajectoryService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteTrajectory(ctx, id)
}
