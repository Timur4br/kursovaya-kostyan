package controllers

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type VideoSegmentService interface {
	Create(context.Context, dto.VideoSegmentInput) (models.VideoSegment, error)
	Get(context.Context, int64) (models.VideoSegment, error)
	Update(context.Context, int64, dto.VideoSegmentInput) (models.VideoSegment, error)
	Delete(context.Context, int64) error
}
type VideoSegmentController struct{ service VideoSegmentService }

func NewVideoSegmentController(service VideoSegmentService) *VideoSegmentController {
	return &VideoSegmentController{service: service}
}
func (c *VideoSegmentController) Create(ctx context.Context, input dto.VideoSegmentInput) (models.VideoSegment, error) {
	return c.service.Create(ctx, input)
}
func (c *VideoSegmentController) Get(ctx context.Context, id int64) (models.VideoSegment, error) {
	return c.service.Get(ctx, id)
}
func (c *VideoSegmentController) Update(ctx context.Context, id int64, input dto.VideoSegmentInput) (models.VideoSegment, error) {
	return c.service.Update(ctx, id, input)
}
func (c *VideoSegmentController) Delete(ctx context.Context, id int64) error {
	return c.service.Delete(ctx, id)
}
