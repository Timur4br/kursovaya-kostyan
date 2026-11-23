package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type TrajectoryPointService interface {
	Create(context.Context, dto.TrajectoryPointInput) (models.TrajectoryPoint, error)
	Get(context.Context, int64) (models.TrajectoryPoint, error)
	Update(context.Context, int64, dto.TrajectoryPointInput) (models.TrajectoryPoint, error)
	Delete(context.Context, int64) error
}
type TrajectoryPointController struct{ service TrajectoryPointService }

func NewTrajectoryPointController(service TrajectoryPointService) *TrajectoryPointController {
	return &TrajectoryPointController{service: service}
}
func (c *TrajectoryPointController) Create(ctx context.Context, input dto.TrajectoryPointInput) (models.TrajectoryPoint, error) {
	return c.service.Create(ctx, input)
}
func (c *TrajectoryPointController) Get(ctx context.Context, id int64) (models.TrajectoryPoint, error) {
	return c.service.Get(ctx, id)
}
func (c *TrajectoryPointController) Update(ctx context.Context, id int64, input dto.TrajectoryPointInput) (models.TrajectoryPoint, error) {
	return c.service.Update(ctx, id, input)
}
func (c *TrajectoryPointController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
