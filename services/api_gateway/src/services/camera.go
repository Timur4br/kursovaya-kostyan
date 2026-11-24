package services

import (
	"log"

	pgDAO "api_gateway/daos/postgres"

	"api_gateway/entities"
)

type CameraService struct {
	dao *pgDAO.CameraDAO
}

func NewCameraService(dao *pgDAO.CameraDAO) *CameraService {
	return &CameraService{
		dao: dao,
	}
}

func (s *CameraService) CreateCamera(camera *entities.Camera) (int, error) {
	id, err := s.dao.CreateCamera(camera)
	if err != nil {
		log.Printf("Error create camera: %v", err)
		return 0, err
	}

	return id, err
}

func (s *CameraService) GetCameraByID(id int) (*entities.Camera, error) {
	camera, err := s.dao.GetCameraByID(id)
	if err != nil {
		log.Printf("Error get camera: %v", err)
		return nil, err
	}

	return camera, err
}

func (s *CameraService) GetCameraList(limit, offset int) ([]entities.Camera, error) {
	cameras, err := s.dao.GetCameraList(limit, offset)
	if err != nil {
		log.Printf("Error get camera list: %v", err)
		return nil, err
	}

	return cameras, err
}

func (s *CameraService) GetCameraCount() (int, error) {
	count, err := s.dao.GetCount()
	if err != nil {
		log.Printf("Error get camera count: %v", err)
		return 0, err
	}

	return count, err
}

func (s *CameraService) UpdateCamera(id int, camera *entities.Camera) error {
	existingCamera, err := s.GetCameraByID(id)
	if err != nil {
		log.Printf("Error get camera for update: %v", err)
		return err
	}

	// TODO:
	// Закинуть в s.dao.UpdateCamera сущность с полями, которые нужно обновить, если они были переданы в s.UpdateCamera
	// Если поля не были переданы, то оставить их без изменений (использовать существующие значения из existingCamera)
	// Возможно добавить структуру для Update методов и использовать её в качестве аргумента для s.UpdateCamera, чтобы явно указать, какие поля нужно обновить
	// Поля в этой структуре могут быть указаны как указатели, чтобы отличать "не передано" от "передано со значением по умолчанию"

	err = s.dao.UpdateCamera(existingCamera)
	if err != nil {
		log.Printf("Error update camera: %v", err)
		return err
	}

	return err
}

func (s *CameraService) DeleteCamera(id int) error {
	err := s.dao.DeleteCamera(id)
	if err != nil {
		log.Printf("Error delete camera: %v", err)
		return err
	}

	return err
}
