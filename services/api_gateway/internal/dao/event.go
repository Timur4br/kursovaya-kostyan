package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateEvent(ctx context.Context, values map[string]any) (models.Event, error) {
	data, err := d.create(ctx, resources.Must("event"), values)
	return decode[models.Event](data, err)
}
func (d *DBDAO) GetEvent(ctx context.Context, id int64) (models.Event, error) {
	data, err := d.get(ctx, resources.Must("event"), []int64{id})
	return decode[models.Event](data, err)
}
func (d *DBDAO) UpdateEvent(ctx context.Context, id int64, values map[string]any) (models.Event, error) {
	data, err := d.replace(ctx, resources.Must("event"), []int64{id}, values)
	return decode[models.Event](data, err)
}
func (d *DBDAO) DeleteEvent(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("event"), []int64{id})
}
