package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateObject(ctx context.Context, values map[string]any) (models.Object, error) {
	data, err := d.create(ctx, resources.Must("object"), values)
	return decode[models.Object](data, err)
}
func (d *DBDAO) GetObject(ctx context.Context, id int64) (models.Object, error) {
	data, err := d.get(ctx, resources.Must("object"), []int64{id})
	return decode[models.Object](data, err)
}
func (d *DBDAO) UpdateObject(ctx context.Context, id int64, values map[string]any) (models.Object, error) {
	data, err := d.replace(ctx, resources.Must("object"), []int64{id}, values)
	return decode[models.Object](data, err)
}
func (d *DBDAO) DeleteObject(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("object"), []int64{id})
}
