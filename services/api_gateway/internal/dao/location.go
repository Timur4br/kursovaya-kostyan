package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateLocation(ctx context.Context, values map[string]any) (models.Location, error) {
	data, err := d.create(ctx, resources.Must("location"), values)
	return decode[models.Location](data, err)
}
func (d *DBDAO) GetLocation(ctx context.Context, id int64) (models.Location, error) {
	data, err := d.get(ctx, resources.Must("location"), []int64{id})
	return decode[models.Location](data, err)
}
func (d *DBDAO) UpdateLocation(ctx context.Context, id int64, values map[string]any) (models.Location, error) {
	data, err := d.replace(ctx, resources.Must("location"), []int64{id}, values)
	return decode[models.Location](data, err)
}
func (d *DBDAO) DeleteLocation(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("location"), []int64{id})
}
