package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateCamera(ctx context.Context, values map[string]any) (models.Camera, error) {
	data, err := d.create(ctx, resources.Must("camera"), values)
	return decode[models.Camera](data, err)
}
func (d *DBDAO) GetCamera(ctx context.Context, id int64) (models.Camera, error) {
	data, err := d.get(ctx, resources.Must("camera"), []int64{id})
	return decode[models.Camera](data, err)
}
func (d *DBDAO) UpdateCamera(ctx context.Context, id int64, values map[string]any) (models.Camera, error) {
	data, err := d.replace(ctx, resources.Must("camera"), []int64{id}, values)
	return decode[models.Camera](data, err)
}
func (d *DBDAO) DeleteCamera(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("camera"), []int64{id})
}
