package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type EventTypeDAO interface {
	CreateEventType(context.Context, map[string]any) (models.EventType, error)
	GetEventType(context.Context, int64) (models.EventType, error)
	UpdateEventType(context.Context, int64, map[string]any) (models.EventType, error)
	DeleteEventType(context.Context, int64) error
}
type EventTypeService struct{ dao EventTypeDAO }

func NewEventTypeService(dao EventTypeDAO) *EventTypeService { return &EventTypeService{dao: dao} }

func (s *EventTypeService) Create(ctx context.Context, input dto.EventTypeInput) (models.EventType, error) {
	values, err := resources.Must("event_type").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.EventType{}, err
	}
	return s.dao.CreateEventType(ctx, values)
}
func (s *EventTypeService) Get(ctx context.Context, id int64) (models.EventType, error) {
	if err := validIDs(id); err != nil {
		return models.EventType{}, err
	}
	return s.dao.GetEventType(ctx, id)
}
func (s *EventTypeService) Update(ctx context.Context, id int64, input dto.EventTypeInput) (models.EventType, error) {
	if err := validIDs(id); err != nil {
		return models.EventType{}, err
	}
	values, err := resources.Must("event_type").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.EventType{}, err
	}
	return s.dao.UpdateEventType(ctx, id, values)
}
func (s *EventTypeService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteEventType(ctx, id)
}
