package services

import (
	"context"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type QueryHistoryDAO interface {
	CreateQueryHistory(context.Context, map[string]any) (models.QueryHistory, error)
	GetQueryHistory(context.Context, int64) (models.QueryHistory, error)
	UpdateQueryHistory(context.Context, int64, map[string]any) (models.QueryHistory, error)
	DeleteQueryHistory(context.Context, int64) error
}
type QueryHistoryService struct{ dao QueryHistoryDAO }

func NewQueryHistoryService(dao QueryHistoryDAO) *QueryHistoryService {
	return &QueryHistoryService{dao: dao}
}

func (s *QueryHistoryService) Create(ctx context.Context, input dto.QueryHistoryInput) (models.QueryHistory, error) {
	values, err := resources.Must("query_history").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.QueryHistory{}, err
	}
	return s.dao.CreateQueryHistory(ctx, values)
}
func (s *QueryHistoryService) Get(ctx context.Context, id int64) (models.QueryHistory, error) {
	if err := validIDs(id); err != nil {
		return models.QueryHistory{}, err
	}
	return s.dao.GetQueryHistory(ctx, id)
}
func (s *QueryHistoryService) Update(ctx context.Context, id int64, input dto.QueryHistoryInput) (models.QueryHistory, error) {
	if err := validIDs(id); err != nil {
		return models.QueryHistory{}, err
	}
	values, err := resources.Must("query_history").ValidateBody(input.Fields(), false)
	if err != nil {
		return models.QueryHistory{}, err
	}
	return s.dao.UpdateQueryHistory(ctx, id, values)
}
func (s *QueryHistoryService) Delete(ctx context.Context, id int64) error {
	if err := validIDs(id); err != nil {
		return err
	}
	return s.dao.DeleteQueryHistory(ctx, id)
}
