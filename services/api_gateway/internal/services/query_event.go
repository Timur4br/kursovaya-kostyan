package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type QueryEventDAO interface {
	CreateQueryEvent(context.Context, map[string]any) (models.QueryEvent, error)
	GetQueryEvent(context.Context, int64, int64) (models.QueryEvent, error)
	UpdateQueryEvent(context.Context, int64, int64, map[string]any) (models.QueryEvent, error)
	DeleteQueryEvent(context.Context, int64, int64) error
}
type QueryEventService struct{ dao QueryEventDAO }

func NewQueryEventService(dao QueryEventDAO) *QueryEventService { return &QueryEventService{dao: dao} }

func (s *QueryEventService) Create(ctx context.Context, input dto.QueryEventInput) (models.QueryEvent, error) {
	values, err := resources.Must("query_event").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.QueryEvent{}, err
	}
	return s.dao.CreateQueryEvent(ctx, values)
}
func (s *QueryEventService) Get(ctx context.Context, queryID int64, eventID int64) (models.QueryEvent, error) {
	if err := validIDs(queryID, eventID); err != nil {
		return models.QueryEvent{}, err
	}
	return s.dao.GetQueryEvent(ctx, queryID, eventID)
}
func (s *QueryEventService) Update(ctx context.Context, queryID int64, eventID int64, input dto.QueryEventInput) (models.QueryEvent, error) {
	if err := validIDs(queryID, eventID); err != nil {
		return models.QueryEvent{}, err
	}
	values, err := resources.Must("query_event").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.QueryEvent{}, err
	}
	return s.dao.UpdateQueryEvent(ctx, queryID, eventID, values)
}
func (s *QueryEventService) Delete(ctx context.Context, queryID int64, eventID int64) error {
	if err := validIDs(queryID, eventID); err != nil {
		return err
	}
	return s.dao.DeleteQueryEvent(ctx, queryID, eventID)
}
