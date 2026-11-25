package services

import (
	"log"

	pgDAO "api_gateway/daos/postgres"

	"api_gateway/entities"
)

type ObjectService struct {
	dao *pgDAO.ObjectDAO
}

func NewObjectService(dao *pgDAO.ObjectDAO) *ObjectService {
	return &ObjectService{
		dao: dao,
	}
}

func (s *ObjectService) CreateObject(object *entities.Object) (int, error) {
	id, err := s.dao.CreateObject(object)
	if err != nil {
		log.Printf("Error create object: %v", err)
		return 0, err
	}

	return id, err
}

func (s *ObjectService) GetObjectByID(id int) (*entities.Object, error) {
	object, err := s.dao.GetObjectByID(id)
	if err != nil {
		log.Printf("Error get object: %v", err)
		return nil, err
	}

	return object, err
}

func (s *ObjectService) GetObjectList(limit, offset int, filter *entities.Filter) ([]entities.Object, error) {
	objects, err := s.dao.GetObjectList(limit, offset, filter)
	if err != nil {
		log.Printf("Error get object list: %v", err)
		return nil, err
	}

	return objects, err
}

func (s *ObjectService) GetObjectCount(filter *entities.Filter) (int, error) {
	count, err := s.dao.GetCount(filter)
	if err != nil {
		log.Printf("Error get object count: %v", err)
		return 0, err
	}

	return count, err
}

func (s *ObjectService) UpdateObject(id int, object *entities.Object) error {
	existingObject, err := s.GetObjectByID(id)
	if err != nil {
		log.Printf("Error get object for update: %v", err)
		return err
	}

	existingObject.Type = object.Type

	err = s.dao.UpdateObject(existingObject)
	if err != nil {
		log.Printf("Error update object: %v", err)
		return err
	}

	return err
}

func (s *ObjectService) DeleteObject(id int) error {
	err := s.dao.DeleteObject(id)
	if err != nil {
		log.Printf("Error delete object: %v", err)
		return err
	}

	return err
}
