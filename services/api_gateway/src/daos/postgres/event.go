package pg_dao

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pgvector/pgvector-go"

	"api_gateway/entities"
)

type EventDAO struct {
	tableName string
	db        *sql.DB
}

func NewEventDAO(db *sql.DB) *EventDAO {
	return &EventDAO{tableName: "events", db: db}
}

func (d *EventDAO) CreateEvent(event *entities.Event) (string, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (camera_id, object_id, type, confidence, embedding, photo_uuid)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING uuid
	`, d.tableName)

	var uuid string
	err := d.db.QueryRow(query, event.CameraID, event.ObjectID, event.Type, event.Confidence, pgvector.NewVector(event.Embedding), event.PhotoUUID).Scan(&uuid)
	if err != nil {
		log.Printf("Error create event query: %v", err)
		return "", err
	}
	return uuid, err
}

func scanEvent(row interface{ Scan(...any) error }) (*entities.Event, error) {
	event := &entities.Event{}
	var embedding pgvector.Vector
	err := row.Scan(&event.UUID, &event.CameraID, &event.ObjectID, &event.Type, &event.Confidence, &embedding, &event.PhotoUUID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	event.Embedding = embedding.Slice()
	return event, nil
}

func (d *EventDAO) GetEventByUUID(uuid string) (*entities.Event, error) {
	query := fmt.Sprintf(`
	SELECT uuid, camera_id, object_id, type, confidence, embedding, photo_uuid, created_at, updated_at
	FROM %s
	WHERE uuid = $1
	`, d.tableName)

	event, err := scanEvent(d.db.QueryRow(query, uuid))
	if err != nil {
		log.Printf("Error get event by [uuid=%v]: %v", uuid, err)
		return nil, err
	}
	return event, err
}

func (d *EventDAO) GetEventList(limit, offset int, filter *entities.Filter) ([]entities.Event, error) {
	query := fmt.Sprintf(`
	SELECT uuid, camera_id, object_id, type, confidence, embedding, photo_uuid, created_at, updated_at
	FROM %s
	`, d.tableName)
	var args []interface{}
	query = addFilterByEvent(query, &args, filter)
	query += " ORDER BY created_at DESC, uuid DESC"
	query = addPagination(query, &args, limit, offset)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		log.Printf("Error get event list: %v", err)
		return nil, err
	}
	defer rows.Close()

	events := make([]entities.Event, 0)
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			log.Printf("Error scan event: %v", err)
			return nil, err
		}
		events = append(events, *event)
	}
	return events, rows.Err()
}

func (d *EventDAO) GetCount(filter *entities.Filter) (int, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, d.tableName)
	var args []any
	query = addFilterByEvent(query, &args, filter)

	var count int
	err := d.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		log.Printf("Error get events count: %v", err)
		return 0, err
	}
	return count, err
}

func (d *EventDAO) UpdateEvent(event *entities.Event) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET camera_id = $2,
	    object_id = $3,
	    type = $4,
	    confidence = $5,
	    embedding = $6,
	    photo_uuid = $7,
	    updated_at = NOW()
	WHERE uuid = $1
	RETURNING uuid
	`, d.tableName)

	var uuid string
	err := d.db.QueryRow(query, event.UUID, event.CameraID, event.ObjectID, event.Type, event.Confidence, pgvector.NewVector(event.Embedding), event.PhotoUUID).Scan(&uuid)
	if err != nil {
		log.Printf("Error update event: %v", err)
		return err
	}
	return err
}

func (d *EventDAO) DeleteEvent(uuid string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE uuid = $1`, d.tableName)
	result, err := d.db.Exec(query, uuid)
	if err != nil {
		log.Printf("Error delete event by [uuid=%v]: %v", uuid, err)
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
