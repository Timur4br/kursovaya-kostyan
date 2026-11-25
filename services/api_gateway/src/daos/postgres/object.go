package pg_dao

import (
	"database/sql"
	"fmt"
	"log"

	"api_gateway/entities"
)

type ObjectDAO struct {
	tableName string
	db        *sql.DB
}

func NewObjectDAO(db *sql.DB) *ObjectDAO {
	return &ObjectDAO{tableName: "objects", db: db}
}

func (d *ObjectDAO) CreateObject(object *entities.Object) (int, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (type)
	VALUES ($1)
	RETURNING id
	`, d.tableName)

	var id int
	err := d.db.QueryRow(query, object.Type).Scan(&id)
	if err != nil {
		log.Printf("Error create object query: %v", err)
		return 0, err
	}
	return id, err
}

func (d *ObjectDAO) GetObjectByID(id int) (*entities.Object, error) {
	query := fmt.Sprintf(`
	SELECT id, type, created_at, updated_at
	FROM %s
	WHERE id = $1
	`, d.tableName)

	object := &entities.Object{}
	err := d.db.QueryRow(query, id).Scan(&object.ID, &object.Type, &object.CreatedAt, &object.UpdatedAt)
	if err != nil {
		log.Printf("Error get object by [id=%v]: %v", id, err)
		return nil, err
	}
	return object, err
}

func (d *ObjectDAO) GetObjectList(limit, offset int, filter *entities.Filter) ([]entities.Object, error) {
	query := fmt.Sprintf(`
	SELECT id, type, created_at, updated_at
	FROM %s
	`, d.tableName)
	var args []interface{}
	query = addFilterByObject(query, &args, filter)
	query += " ORDER BY id"
	query = addPagination(query, &args, limit, offset)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		log.Printf("Error get object list: %v", err)
		return nil, err
	}
	defer rows.Close()

	objects := make([]entities.Object, 0)
	for rows.Next() {
		object := entities.Object{}
		err := rows.Scan(&object.ID, &object.Type, &object.CreatedAt, &object.UpdatedAt)
		if err != nil {
			log.Printf("Error scan object: %v", err)
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, rows.Err()
}

func (d *ObjectDAO) GetCount(filter *entities.Filter) (int, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, d.tableName)
	var args []any
	query = addFilterByObject(query, &args, filter)

	var count int
	err := d.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		log.Printf("Error get objects count: %v", err)
		return 0, err
	}
	return count, err
}

func (d *ObjectDAO) UpdateObject(object *entities.Object) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET type = $2,
	    updated_at = NOW()
	WHERE id = $1
	RETURNING id
	`, d.tableName)

	var id int
	err := d.db.QueryRow(query, object.ID, object.Type).Scan(&id)
	if err != nil {
		log.Printf("Error update object: %v", err)
		return err
	}
	return err
}

func (d *ObjectDAO) DeleteObject(id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, d.tableName)
	result, err := d.db.Exec(query, id)
	if err != nil {
		log.Printf("Error delete object by [id=%v]: %v", id, err)
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
