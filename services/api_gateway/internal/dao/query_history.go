package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreateQueryHistory(ctx context.Context, values map[string]any) (models.QueryHistory, error) {
	data, err := d.create(ctx, resources.Must("query_history"), values)
	return decode[models.QueryHistory](data, err)
}
func (d *DBDAO) GetQueryHistory(ctx context.Context, id int64) (models.QueryHistory, error) {
	data, err := d.get(ctx, resources.Must("query_history"), []int64{id})
	return decode[models.QueryHistory](data, err)
}
func (d *DBDAO) UpdateQueryHistory(ctx context.Context, id int64, values map[string]any) (models.QueryHistory, error) {
	data, err := d.replace(ctx, resources.Must("query_history"), []int64{id}, values)
	return decode[models.QueryHistory](data, err)
}
func (d *DBDAO) DeleteQueryHistory(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("query_history"), []int64{id})
}
