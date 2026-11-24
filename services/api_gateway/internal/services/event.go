package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type EventDAO interface {
	CreateEvent(context.Context, map[string]any) (models.Event, error)
	GetEvent(context.Context, int64) (models.Event, error)
	UpdateEvent(context.Context, int64, map[string]any) (models.Event, error)
	DeleteEvent(context.Context, int64) error
}
type EventService struct{ dao EventDAO }

func NewEventService(dao EventDAO) *EventService { return &EventService{dao: dao} }

func (s *EventService) Create(ctx context.Context, input dto.EventInput) (models.Event, error) {
	values, err := resources.Must("event").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Event{}, err
	}
	return s.dao.CreateEvent(ctx, values)
}
func (s *EventService) Get(ctx context.Context, id int64) (models.Event, error) {
	if err := validIDs(id); err != nil {
		return models.Event{}, err
	}
	return s.dao.GetEvent(ctx, id)
}
func (s *EventService) Update(ctx context.Context, id int64, input dto.EventInput) (models.Event, error) {
	if err := validIDs(id); err != nil {
		return models.Event{}, err
	}
	values, err := resources.Must("event").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.Event{}, err
	}
	return s.dao.UpdateEvent(ctx, id, values)
}
func (s *EventService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteEvent(ctx, id)
}
