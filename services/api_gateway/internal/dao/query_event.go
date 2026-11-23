package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateQueryEvent(ctx context.Context, values map[string]any) (models.QueryEvent, error) {
	data, err := d.create(ctx, resources.Must("query_event"), values)
	return decode[models.QueryEvent](data, err)
}
func (d *DBDAO) GetQueryEvent(ctx context.Context, queryID int64, eventID int64) (models.QueryEvent, error) {
	data, err := d.get(ctx, resources.Must("query_event"), []int64{queryID, eventID})
	return decode[models.QueryEvent](data, err)
}
func (d *DBDAO) UpdateQueryEvent(ctx context.Context, queryID int64, eventID int64, values map[string]any) (models.QueryEvent, error) {
	data, err := d.replace(ctx, resources.Must("query_event"), []int64{queryID, eventID}, values)
	return decode[models.QueryEvent](data, err)
}
func (d *DBDAO) DeleteQueryEvent(ctx context.Context, queryID int64, eventID int64) error {
	return d.delete(ctx, resources.Must("query_event"), []int64{queryID, eventID})
}
