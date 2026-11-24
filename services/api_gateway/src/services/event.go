package services

import (
	"log"

	pgDAO "api_gateway/daos/postgres"

	"api_gateway/entities"
)

type EventService struct {
	dao *pgDAO.EventDAO
}

func NewEventService(dao *pgDAO.EventDAO) *EventService {
	return &EventService{
		dao: dao,
	}
}

func (s *EventService) CreateEvent(event *entities.Event) (string, error) {
	uuid, err := s.dao.CreateEvent(event)
	if err != nil {
		log.Printf("Error create event: %v", err)
		return "", err
	}

	return uuid, err
}

func (s *EventService) GetEventByUUID(uuid string) (*entities.Event, error) {
	event, err := s.dao.GetEventByUUID(uuid)
	if err != nil {
		log.Printf("Error get event: %v", err)
		return nil, err
	}

	return event, err
}

func (s *EventService) GetEventList(limit, offset int, filter *entities.Filter) ([]entities.Event, error) {
	events, err := s.dao.GetEventList(limit, offset, filter)
	if err != nil {
		log.Printf("Error get event list: %v", err)
		return nil, err
	}

	return events, err
}

func (s *EventService) GetEventCount(filter *entities.Filter) (int, error) {
	count, err := s.dao.GetCount(filter)
	if err != nil {
		log.Printf("Error get event count: %v", err)
		return 0, err
	}

	return count, err
}

func (s *EventService) UpdateEvent(uuid string, event *entities.Event) error {
	existingEvent, err := s.GetEventByUUID(uuid)
	if err != nil {
		log.Printf("Error get event for update: %v", err)
		return err
	}

	existingEvent.CameraID = event.CameraID
	existingEvent.ObjectID = event.ObjectID
	existingEvent.Type = event.Type
	existingEvent.Confidence = event.Confidence
	existingEvent.Embedding = event.Embedding
	existingEvent.PhotoUUID = event.PhotoUUID

	err = s.dao.UpdateEvent(existingEvent)
	if err != nil {
		log.Printf("Error update event: %v", err)
		return err
	}

	return err
}

func (s *EventService) DeleteEvent(uuid string) error {
	err := s.dao.DeleteEvent(uuid)
	if err != nil {
		log.Printf("Error delete event: %v", err)
		return err
	}

	return err
}
