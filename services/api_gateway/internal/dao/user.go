package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateUser(ctx context.Context, values map[string]any) (models.User, error) {
	data, err := d.create(ctx, resources.Must("user"), values)
	return decode[models.User](data, err)
}
func (d *DBDAO) GetUser(ctx context.Context, id int64) (models.User, error) {
	data, err := d.get(ctx, resources.Must("user"), []int64{id})
	return decode[models.User](data, err)
}
func (d *DBDAO) UpdateUser(ctx context.Context, id int64, values map[string]any) (models.User, error) {
	data, err := d.replace(ctx, resources.Must("user"), []int64{id}, values)
	return decode[models.User](data, err)
}
func (d *DBDAO) DeleteUser(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("user"), []int64{id})
}
