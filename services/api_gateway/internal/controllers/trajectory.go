package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type TrajectoryService interface {
	Create(context.Context, dto.TrajectoryInput) (models.Trajectory, error)
	Get(context.Context, int64) (models.Trajectory, error)
	Update(context.Context, int64, dto.TrajectoryInput) (models.Trajectory, error)
	Delete(context.Context, int64) error
}
type TrajectoryController struct{ service TrajectoryService }

func NewTrajectoryController(service TrajectoryService) *TrajectoryController {
	return &TrajectoryController{service: service}
}
func (c *TrajectoryController) Create(ctx context.Context, input dto.TrajectoryInput) (models.Trajectory, error) {
	return c.service.Create(ctx, input)
}
func (c *TrajectoryController) Get(ctx context.Context, id int64) (models.Trajectory, error) {
	return c.service.Get(ctx, id)
}
func (c *TrajectoryController) Update(ctx context.Context, id int64, input dto.TrajectoryInput) (models.Trajectory, error) {
	return c.service.Update(ctx, id, input)
}
func (c *TrajectoryController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
