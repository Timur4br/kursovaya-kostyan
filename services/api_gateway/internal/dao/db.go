// Package dao implements entity persistence using parameterized PostgreSQL queries.
package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/resources"
)

type DBDAO struct{ db *sql.DB }

func New(db *sql.DB) *DBDAO { return &DBDAO{db: db} }

func quote(identifier string) string { return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"` }
func projection(r resources.Resource) string {
	parts := make([]string, 0, len(r.Fields)*2)
	for _, f := range r.Fields {
		parts = append(parts, "'"+f.Name+"'", quote(f.Name))
	}
	return "jsonb_build_object(" + strings.Join(parts, ",") + ")"
}
func keyWhere(r resources.Resource, ids []int64, offset int) (string, []any) {
	parts := make([]string, len(r.Keys))
	values := make([]any, len(ids))
	for i, key := range r.Keys {
		parts[i] = quote(key) + "=$" + strconv.Itoa(offset+i+1)
		values[i] = ids[i]
	}
	return strings.Join(parts, " AND "), values
}
func decode[T any](data json.RawMessage, err error) (T, error) {
	var record T
	if err != nil {
		return record, err
	}
	err = json.Unmarshal(data, &record)
	return record, err
}

func (d *DBDAO) get(ctx context.Context, r resources.Resource, ids []int64) (json.RawMessage, error) {
	where, values := keyWhere(r, ids, 0)
	var data json.RawMessage
	err := d.db.QueryRowContext(ctx, "SELECT "+projection(r)+" FROM "+quote(r.Table)+" WHERE "+where, values...).Scan(&data)
	return data, err
}

func (d *DBDAO) create(ctx context.Context, r resources.Resource, values map[string]any) (json.RawMessage, error) {
	columns, params := []string{}, []string{}
	args := []any{}
	for _, f := range r.Fields {
		if value, ok := values[f.Name]; ok && !f.ReadOnly {
			columns = append(columns, quote(f.Name))
			args = append(args, value)
			params = append(params, "$"+strconv.Itoa(len(args)))
		}
	}
	query := "INSERT INTO " + quote(r.Table)
	if len(columns) == 0 {
		query += " DEFAULT VALUES"
	} else {
		query += " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Join(params, ",") + ")"
	}
	query += " RETURNING " + projection(r)
	var data json.RawMessage
	err := d.db.QueryRowContext(ctx, query, args...).Scan(&data)
	return data, err
}

// replace implements PUT: omitted optional fields are reset to database defaults.
// Immutable fields such as id, created_at and generated duration are preserved.
func (d *DBDAO) replace(ctx context.Context, r resources.Resource, ids []int64, values map[string]any) (json.RawMessage, error) {
	sets := []string{}
	args := []any{}
	for _, f := range r.Fields {
		if f.ReadOnly {
			continue
		}
		if value, ok := values[f.Name]; ok {
			args = append(args, value)
			sets = append(sets, quote(f.Name)+"=$"+strconv.Itoa(len(args)))
		} else {
			sets = append(sets, quote(f.Name)+"=DEFAULT")
		}
	}
	where, keyArgs := keyWhere(r, ids, len(args))
	args = append(args, keyArgs...)
	var data json.RawMessage
	err := d.db.QueryRowContext(ctx, "UPDATE "+quote(r.Table)+" SET "+strings.Join(sets, ",")+" WHERE "+where+" RETURNING "+projection(r), args...).Scan(&data)
	return data, err
}

func (d *DBDAO) delete(ctx context.Context, r resources.Resource, ids []int64) error {
	where, args := keyWhere(r, ids, 0)
	result, err := d.db.ExecContext(ctx, "DELETE FROM "+quote(r.Table)+" WHERE "+where, args...)
	if err != nil {
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
