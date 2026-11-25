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

func (s *CameraService) GetCameraList(limit, offset int, filter *entities.Filter) ([]entities.Camera, error) {
	cameras, err := s.dao.GetCameraList(limit, offset, filter)
	if err != nil {
		log.Printf("Error get camera list: %v", err)
		return nil, err
	}

	return cameras, err
}

func (s *CameraService) GetCameraCount(filter *entities.Filter) (int, error) {
	count, err := s.dao.GetCount(filter)
	if err != nil {
		log.Printf("Error get camera count: %v", err)
		return 0, err
	}

	return count, err
}

func (s *CameraService) UpdateCamera(id int, camera *entities.CameraUpdate) error {
	existingCamera, err := s.GetCameraByID(id)
	if err != nil {
		log.Printf("Error get camera for update: %v", err)
		return err
	}

	if camera.Comment != nil {
		existingCamera.Comment = *camera.Comment
	}
	if camera.IP != nil {
		existingCamera.IP = *camera.IP
	}
	if camera.Port != nil {
		existingCamera.Port = *camera.Port
	}
	if camera.StreamURL != nil {
		existingCamera.StreamURL = *camera.StreamURL
	}
	if camera.IsEnabled != nil {
		existingCamera.IsEnabled = *camera.IsEnabled
	}

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
