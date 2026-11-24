package dao

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

func (d *DBDAO) CreatePerson(ctx context.Context, values map[string]any) (models.Person, error) {
	data, err := d.create(ctx, resources.Must("person"), values)
	return decode[models.Person](data, err)
}
func (d *DBDAO) GetPerson(ctx context.Context, id int64) (models.Person, error) {
	data, err := d.get(ctx, resources.Must("person"), []int64{id})
	return decode[models.Person](data, err)
}
func (d *DBDAO) UpdatePerson(ctx context.Context, id int64, values map[string]any) (models.Person, error) {
	data, err := d.replace(ctx, resources.Must("person"), []int64{id}, values)
	return decode[models.Person](data, err)
}
func (d *DBDAO) DeletePerson(ctx context.Context, id int64) error {
	return d.delete(ctx, resources.Must("person"), []int64{id})
}
