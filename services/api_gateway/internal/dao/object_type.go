package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateObjectType(ctx context.Context, values map[string]any) (models.ObjectType, error) {
	data, err := d.create(ctx, resources.Must("object_type"), values)
	return decode[models.ObjectType](data, err)
}
func (d *DBDAO) GetObjectType(ctx context.Context, id int64) (models.ObjectType, error) {
	data, err := d.get(ctx, resources.Must("object_type"), []int64{id})
	return decode[models.ObjectType](data, err)
}
func (d *DBDAO) UpdateObjectType(ctx context.Context, id int64, values map[string]any) (models.ObjectType, error) {
	data, err := d.replace(ctx, resources.Must("object_type"), []int64{id}, values)
	return decode[models.ObjectType](data, err)
}
func (d *DBDAO) DeleteObjectType(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("object_type"), []int64{id})
}
