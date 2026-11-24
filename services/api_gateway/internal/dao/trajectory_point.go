package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateTrajectoryPoint(ctx context.Context, values map[string]any) (models.TrajectoryPoint, error) {
	data, err := d.create(ctx, resources.Must("trajectory_point"), values)
	return decode[models.TrajectoryPoint](data, err)
}
func (d *DBDAO) GetTrajectoryPoint(ctx context.Context, id int64) (models.TrajectoryPoint, error) {
	data, err := d.get(ctx, resources.Must("trajectory_point"), []int64{id})
	return decode[models.TrajectoryPoint](data, err)
}
func (d *DBDAO) UpdateTrajectoryPoint(ctx context.Context, id int64, values map[string]any) (models.TrajectoryPoint, error) {
	data, err := d.replace(ctx, resources.Must("trajectory_point"), []int64{id}, values)
	return decode[models.TrajectoryPoint](data, err)
}
func (d *DBDAO) DeleteTrajectoryPoint(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("trajectory_point"), []int64{id})
}
