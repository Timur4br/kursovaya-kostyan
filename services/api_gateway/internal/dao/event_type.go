package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateEventType(ctx context.Context, values map[string]any) (models.EventType, error) {
	data, err := d.create(ctx, resources.Must("event_type"), values)
	return decode[models.EventType](data, err)
}
func (d *DBDAO) GetEventType(ctx context.Context, id int64) (models.EventType, error) {
	data, err := d.get(ctx, resources.Must("event_type"), []int64{id})
	return decode[models.EventType](data, err)
}
func (d *DBDAO) UpdateEventType(ctx context.Context, id int64, values map[string]any) (models.EventType, error) {
	data, err := d.replace(ctx, resources.Must("event_type"), []int64{id}, values)
	return decode[models.EventType](data, err)
}
func (d *DBDAO) DeleteEventType(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("event_type"), []int64{id})
}
