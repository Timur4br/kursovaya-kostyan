package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateTrajectory(ctx context.Context, values map[string]any) (models.Trajectory, error) {
	data, err := d.create(ctx, resources.Must("trajectory"), values)
	return decode[models.Trajectory](data, err)
}
func (d *DBDAO) GetTrajectory(ctx context.Context, id int64) (models.Trajectory, error) {
	data, err := d.get(ctx, resources.Must("trajectory"), []int64{id})
	return decode[models.Trajectory](data, err)
}
func (d *DBDAO) UpdateTrajectory(ctx context.Context, id int64, values map[string]any) (models.Trajectory, error) {
	data, err := d.replace(ctx, resources.Must("trajectory"), []int64{id}, values)
	return decode[models.Trajectory](data, err)
}
func (d *DBDAO) DeleteTrajectory(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("trajectory"), []int64{id})
}
