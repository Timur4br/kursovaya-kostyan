package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateVideoSegment(ctx context.Context, values map[string]any) (models.VideoSegment, error) {
	data, err := d.create(ctx, resources.Must("video_segment"), values)
	return decode[models.VideoSegment](data, err)
}
func (d *DBDAO) GetVideoSegment(ctx context.Context, id int64) (models.VideoSegment, error) {
	data, err := d.get(ctx, resources.Must("video_segment"), []int64{id})
	return decode[models.VideoSegment](data, err)
}
func (d *DBDAO) UpdateVideoSegment(ctx context.Context, id int64, values map[string]any) (models.VideoSegment, error) {
	data, err := d.replace(ctx, resources.Must("video_segment"), []int64{id}, values)
	return decode[models.VideoSegment](data, err)
}
func (d *DBDAO) DeleteVideoSegment(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("video_segment"), []int64{id})
}
