package entities

import "time"

type Event struct {
	UUID       string    `json:"uuid"`
	CameraID   int       `json:"camera_id" binding:"required,gt=0"`
	ObjectID   int       `json:"object_id" binding:"required,gt=0"`
	Type       string    `json:"type" binding:"required"`
	Confidence float64   `json:"confidence"`
	Embedding  []float32 `json:"embedding" binding:"required,len=512"`
	PhotoUUID  string    `json:"photo_uuid" binding:"required,uuid"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
