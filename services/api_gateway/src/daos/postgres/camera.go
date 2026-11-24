package pg_dao

import (
	"database/sql"
	"fmt"
	"log"

	"api_gateway/entities"
)

type CameraDAO struct {
	tableName string

	db *sql.DB
}

func NewCameraDAO(db *sql.DB) *CameraDAO {
	tableName := "camera"

	return &CameraDAO{
		tableName: tableName,

		db: db,
	}
}

func (d *CameraDAO) CreateCamera(camera *entities.Camera) (int, error) {
	query := fmt.Sprintf(`
	INSERT
		INTO %s
		(comment, stream_url, ip, port, is_enabled)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING id
	`, d.tableName)

	var id int
	err := d.db.QueryRow(query,
		camera.Comment, camera.StreamURL, camera.IP, camera.Port, camera.IsEnabled,
	).Scan(&id)
	if err != nil {
		log.Printf("Error create camera query: %v", err)
		return 0, err
	}

	return id, err
}

func (d *CameraDAO) GetCameraByID(id int) (*entities.Camera, error) {
	query := fmt.Sprintf(`
	SELECT
		id, comment, stream_url, ip, port, is_enabled, created_at, updated_at
	FROM
		%s AS t
	WHERE
		t.id = $1
	`, d.tableName)

	camera := &entities.Camera{}
	err := d.db.QueryRow(query, id).Scan(
		&camera.ID, &camera.Comment, &camera.StreamURL, &camera.IP, &camera.Port, &camera.IsEnabled, &camera.CreatedAt, &camera.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error get camera by [id=%d]: %v", id, err)
		return nil, err
	}

	return camera, err
}

func (d *CameraDAO) GetCameraList(limit, offset int) ([]entities.Camera, error) {
	query := fmt.Sprintf(`
	SELECT
		id, comment, stream_url, ip, port, is_enabled, created_at, updated_at
	FROM
		%s AS c
	`, d.tableName)

	var args []interface{}

	query += " ORDER BY id"
	query = addPagination(query, &args, limit, offset)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		log.Printf("Error get camera list: %v", err)
		return nil, err
	}

	cameras := make([]entities.Camera, 0)

	for rows.Next() {
		camera := entities.Camera{}

		err = rows.Scan(
			&camera.ID, &camera.Comment, &camera.StreamURL, &camera.IP, &camera.Port, &camera.IsEnabled, &camera.CreatedAt, &camera.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scan camera: %v", err)
			return nil, err
		}

		cameras = append(cameras, camera)
	}

	return cameras, err
}

func (d *CameraDAO) GetCount() (int, error) {
	query := fmt.Sprintf(`
	SELECT
		COUNT(*)
	FROM
		%s
	`, d.tableName)

	var count int
	err := d.db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error get cameras count: %v", err)
		return 0, err
	}

	return count, err
}

func (d *CameraDAO) UpdateCamera(camera *entities.Camera) error {
	query := fmt.Sprintf(`
	UPDATE
		%s
	SET
		comment = $2,
		stream_url = $3,
		ip = $4,
		port = $5,
		is_enabled = $6,
		updated_at = NOW()
	WHERE
		id = $1
	`, d.tableName)

	var id int
	err := d.db.QueryRow(query,
		camera.ID,
		camera.Comment, camera.StreamURL, camera.IP, camera.Port, camera.IsEnabled,
	).Scan(&id)
	if err != nil {
		log.Printf("Error update camera: %v", err)
		return err
	}

	return err
}

func (d *CameraDAO) DeleteCamera(id int) error {
	query := fmt.Sprintf(`
	DELETE
	FROM
		%s AS t
	WHERE
		t.id = $1
	`, d.tableName)

	_, err := d.db.Exec(query, id)
	if err != nil {
		log.Printf("Error delete camera by [id=%d]: %v", id, err)
		return err
	}

	return err
}
