package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type CameraService interface {
	Create(context.Context, dto.CameraInput) (models.Camera, error)
	Get(context.Context, int64) (models.Camera, error)
	Update(context.Context, int64, dto.CameraInput) (models.Camera, error)
	Delete(context.Context, int64) error
}
type CameraController struct{ service CameraService }

func NewCameraController(service CameraService) *CameraController {
	return &CameraController{service: service}
}
func (c *CameraController) Create(ctx context.Context, input dto.CameraInput) (models.Camera, error) {
	return c.service.Create(ctx, input)
}
func (c *CameraController) Get(ctx context.Context, id int64) (models.Camera, error) {
	return c.service.Get(ctx, id)
}
func (c *CameraController) Update(ctx context.Context, id int64, input dto.CameraInput) (models.Camera, error) {
	return c.service.Update(ctx, id, input)
}
func (c *CameraController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
