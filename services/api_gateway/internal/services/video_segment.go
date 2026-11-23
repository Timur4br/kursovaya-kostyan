package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type VideoSegmentDAO interface {
	CreateVideoSegment(context.Context, map[string]any) (models.VideoSegment, error)
	GetVideoSegment(context.Context, int64) (models.VideoSegment, error)
	UpdateVideoSegment(context.Context, int64, map[string]any) (models.VideoSegment, error)
	DeleteVideoSegment(context.Context, int64) error
}
type VideoSegmentService struct{ dao VideoSegmentDAO }

func NewVideoSegmentService(dao VideoSegmentDAO) *VideoSegmentService {
	return &VideoSegmentService{dao: dao}
}

func (s *VideoSegmentService) Create(ctx context.Context, input dto.VideoSegmentInput) (models.VideoSegment, error) {
	values, err := resources.Must("video_segment").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.VideoSegment{}, err
	}
	return s.dao.CreateVideoSegment(ctx, values)
}
func (s *VideoSegmentService) Get(ctx context.Context, id int64) (models.VideoSegment, error) {
	if err := validIDs(id); err != nil {
		return models.VideoSegment{}, err
	}
	return s.dao.GetVideoSegment(ctx, id)
}
func (s *VideoSegmentService) Update(ctx context.Context, id int64, input dto.VideoSegmentInput) (models.VideoSegment, error) {
	if err := validIDs(id); err != nil {
		return models.VideoSegment{}, err
	}
	values, err := resources.Must("video_segment").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.VideoSegment{}, err
	}
	return s.dao.UpdateVideoSegment(ctx, id, values)
}
func (s *VideoSegmentService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteVideoSegment(ctx, id)
}
