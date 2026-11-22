package database

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
)

func integrationDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 against PostgreSQL + pgvector")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	admin, err := Open(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("integration_%d", time.Now().UnixNano())
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err = admin.ExecContext(ctx, "CREATE SCHEMA "+quoted); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	connection, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	// Avoid resolving the application's public goose_db_version table in isolated tests.
	connection.RuntimeParams["search_path"] = schema
	db := stdlib.OpenDB(*connection)
	t.Cleanup(func() {
		db.Close()
		// schema is generated above, never supplied by a user or environment variable.
		if _, err := admin.ExecContext(ctx, "DROP SCHEMA "+quoted+" CASCADE"); err != nil {
			t.Error(err)
		}
		admin.Close()
	})
	return db, schema
}

func TestPostgresMigrations(t *testing.T) {
	db, schema := integrationDB(t)
	ctx := context.Background()
	if err := Ready(ctx, db); err == nil {
		t.Fatal("unmigrated database reported ready")
	}
	if err := Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := Ready(ctx, db); err != nil {
		t.Fatal(err)
	}
	var tables int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema=$1 AND table_name<>'goose_db_version'", schema).Scan(&tables); err != nil {
		t.Fatal(err)
	}
	if tables != 13 {
		t.Fatalf("want 13 tables, got %d", tables)
	}
	var distance float64
	if err := db.QueryRowContext(ctx, "SELECT '[1,0]'::public.vector OPERATOR(public.<->) '[0,1]'::public.vector").Scan(&distance); err != nil {
		t.Fatal(err)
	}
	if math.Abs(distance-math.Sqrt2) > 1e-6 {
		t.Fatal(distance)
	}
	if _, err := db.ExecContext(ctx, "INSERT INTO users(username,name) VALUES('migration_probe','Keep this row')"); err != nil {
		t.Fatal(err)
	}
	if err := Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM users WHERE username='migration_probe'").Scan(&count); err != nil || count != 1 {
		t.Fatal("repeated migration lost data", err)
	}
	provider, err := Provider(db)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = provider.Down(ctx); err != nil {
		t.Fatal(err)
	}
	if err := Ready(ctx, db); err == nil {
		t.Fatal("rolled back schema reported ready")
	}
	if err := Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := Ready(ctx, db); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresRelationshipsAndConstraints(t *testing.T) {
	db, _ := integrationDB(t)
	ctx := context.Background()
	if err := Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	id := func(query string, args ...any) int64 {
		t.Helper()
		var value int64
		if err := tx.QueryRowContext(ctx, query, args...).Scan(&value); err != nil {
			t.Fatal(err)
		}
		return value
	}
	user := id("INSERT INTO users(username,name) VALUES('owner','Owner') RETURNING id")
	place := id("INSERT INTO locations(name,type) VALUES('Room','office') RETURNING id")
	otherPlace := id("INSERT INTO locations(name,type) VALUES('Other room','office') RETURNING id")
	cam := id("INSERT INTO cameras(location_id,name,fps) VALUES($1,'Camera',25) RETURNING id", place)
	otherCam := id("INSERT INTO cameras(location_id,name) VALUES($1,'Other camera') RETURNING id", place)
	person := id("INSERT INTO persons(user_id) VALUES($1) RETURNING id", user)
	objectType := id("INSERT INTO object_types(name) VALUES('Backpack') RETURNING id")
	object := id("INSERT INTO objects(object_type_id,name) VALUES($1,'My backpack') RETURNING id", objectType)
	eventType := id("INSERT INTO event_types(name) VALUES('placed') RETURNING id")
	event := id("INSERT INTO events(event_type_id,camera_id,location_id,person_id,object_id,occurred_at,confidence) VALUES($1,$2,$3,$4,$5,now(),0.95) RETURNING id", eventType, cam, place, person, object)
	otherEvent := id("INSERT INTO events(event_type_id,camera_id,location_id,occurred_at) VALUES($1,$2,$3,now()) RETURNING id", eventType, cam, place)
	video := id("INSERT INTO video_segments(camera_id,event_id,start_time,end_time,storage_reference) VALUES($1,$2,now(),now()+interval '12 seconds','test.mp4') RETURNING id", cam, event)
	trajectory := id("INSERT INTO trajectories(person_id,start_time,start_location_id) VALUES($1,now(),$2) RETURNING id", person, place)
	id("INSERT INTO trajectory_points(trajectory_id,video_segment_id,observed_at,x,y) VALUES($1,$2,now(),1,2) RETURNING id", trajectory, video)
	query := id("INSERT INTO query_history(user_id,text) VALUES($1,'Where is my backpack?') RETURNING id", user)
	otherQuery := id("INSERT INTO query_history(user_id,text) VALUES($1,'What happened?') RETURNING id", user)
	for _, q := range []int64{query, otherQuery} {
		for _, e := range []int64{event, otherEvent} {
			if _, err := tx.ExecContext(ctx, "INSERT INTO query_events(query_id,event_id) VALUES($1,$2)", q, e); err != nil {
				t.Fatal(err)
			}
		}
	}
	var duration float64
	if err := tx.QueryRowContext(ctx, "SELECT duration_seconds FROM video_segments WHERE id=$1", video).Scan(&duration); err != nil || duration != 12 {
		t.Fatal("invalid generated duration", duration, err)
	}
	cases := []struct {
		name, sql, state string
		args             []any
	}{
		{"orphan_camera", "INSERT INTO cameras(location_id,name) VALUES(-1,'bad')", "23503", nil},
		{"invalid_fps", "UPDATE cameras SET fps=0 WHERE id=$1", "23514", []any{cam}},
		{"invalid_confidence", "UPDATE events SET confidence=1.1 WHERE id=$1", "23514", []any{event}},
		{"camera_location_mismatch", "UPDATE events SET location_id=$1 WHERE id=$2", "23503", []any{otherPlace, event}},
		{"clip_camera_mismatch", "UPDATE video_segments SET camera_id=$1 WHERE id=$2", "23503", []any{otherCam, video}},
		{"negative_duration", "UPDATE video_segments SET end_time=start_time-interval '1 second' WHERE id=$1", "23514", []any{video}},
		{"inverted_person_time", "UPDATE persons SET first_seen=now(),last_seen=now()-interval '1 second' WHERE id=$1", "23514", []any{person}},
		{"inverted_object_time", "UPDATE objects SET first_seen=now(),last_seen=now()-interval '1 second' WHERE id=$1", "23514", []any{object}},
		{"inverted_trajectory_time", "UPDATE trajectories SET end_time=start_time-interval '1 second' WHERE id=$1", "23514", []any{trajectory}},
		{"invalid_coordinate", "UPDATE trajectory_points SET x='NaN'::FLOAT8 WHERE trajectory_id=$1", "23514", []any{trajectory}},
		{"blank_query", "UPDATE query_history SET text='  ' WHERE id=$1", "23514", []any{query}},
		{"unknown_query_status", "UPDATE query_history SET status='anything' WHERE id=$1", "23514", []any{query}},
		{"duplicate_result_link", "INSERT INTO query_events(query_id,event_id) VALUES($1,$2)", "23505", []any{query, event}},
		{"delete_referenced_camera", "DELETE FROM cameras WHERE id=$1", "23503", []any{cam}},
		{"delete_referenced_user", "DELETE FROM users WHERE id=$1", "23503", []any{user}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if _, err := tx.ExecContext(ctx, "SAVEPOINT invalid_input"); err != nil {
				t.Fatal(err)
			}
			_, err := tx.ExecContext(ctx, test.sql, test.args...)
			pgErr, ok := err.(*pgconn.PgError)
			if !ok || pgErr.Code != test.state {
				t.Errorf("want SQLSTATE %s, got %v", test.state, err)
			}
			if _, err := tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT invalid_input"); err != nil {
				t.Fatal(err)
			}
			if _, err := tx.ExecContext(ctx, "RELEASE SAVEPOINT invalid_input"); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// These two tests are invoked separately around a Compose restart by verification scripts.
func persistenceConnection(t *testing.T) (*sql.DB, string) {
	t.Helper()
	token := os.Getenv("PERSISTENCE_PROBE")
	if token == "" {
		t.Skip("no persistence probe requested")
	}
	if len(token) > 60 {
		t.Fatal("probe token too long")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	db, err := Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db, "__verify_" + token
}

func TestPersistenceWrite(t *testing.T) {
	db, key := persistenceConnection(t)
	if _, err := db.Exec("INSERT INTO users(username,name) VALUES($1,'volume verification')", key); err != nil {
		t.Fatal(err)
	}
}

func TestPersistenceRead(t *testing.T) {
	db, key := persistenceConnection(t)
	var name string
	if err := db.QueryRow("SELECT name FROM users WHERE username=$1", key).Scan(&name); err != nil {
		t.Fatal("row missing after container restart", err)
	}
	if name != "volume verification" {
		t.Fatal("row changed after restart")
	}
	if _, err := db.Exec("DELETE FROM users WHERE username=$1", key); err != nil {
		t.Fatal(err)
	}
}
