package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type PersonDAO interface {
	CreatePerson(context.Context, map[string]any) (models.Person, error)
	GetPerson(context.Context, int64) (models.Person, error)
	UpdatePerson(context.Context, int64, map[string]any) (models.Person, error)
	DeletePerson(context.Context, int64) error
}
type PersonService struct{ dao PersonDAO }

func NewPersonService(dao PersonDAO) *PersonService { return &PersonService{dao: dao} }

func (s *PersonService) Create(ctx context.Context, input dto.PersonInput) (models.Person, error) {
	values, err := resources.Must("person").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Person{}, err
	}
	return s.dao.CreatePerson(ctx, values)
}
func (s *PersonService) Get(ctx context.Context, id int64) (models.Person, error) {
	if err := validIDs(id); err != nil {
		return models.Person{}, err
	}
	return s.dao.GetPerson(ctx, id)
}
func (s *PersonService) Update(ctx context.Context, id int64, input dto.PersonInput) (models.Person, error) {
	if err := validIDs(id); err != nil {
		return models.Person{}, err
	}
	values, err := resources.Must("person").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Person{}, err
	}
	return s.dao.UpdatePerson(ctx, id, values)
}
func (s *PersonService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeletePerson(ctx, id)
}
