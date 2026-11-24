package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Timur4br/kursovaya-kostyan/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

func Open(ctx context.Context, connectionURL string) (*sql.DB, error) {
	cfg, err := pgx.ParseConfig(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database configuration")
	}
	cfg.RuntimeParams["application_name"] = "api_gateway"
	db := stdlib.OpenDB(*cfg)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(4)
	db.SetConnMaxLifetime(30 * time.Minute)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database connection failed")
	}
	return db, nil
}

func Provider(db *sql.DB) (*goose.Provider, error) {
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return nil, err
	}
	return goose.NewProvider(goose.DialectPostgres, db, migrations.Files, goose.WithSessionLocker(locker), goose.WithDisableGlobalRegistry(true))
}

func Migrate(ctx context.Context, db *sql.DB) error {
	provider, err := Provider(db)
	if err != nil {
		return err
	}
	_, err = provider.Up(ctx)
	return err
}

// Ready checks connectivity, the installed extension and the applied schema version.
// It never creates tables or runs migrations from an HTTP request.
func Ready(ctx context.Context, db *sql.DB) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database unavailable")
	}
	var enabled bool
	if err := db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname='vector')").Scan(&enabled); err != nil || !enabled {
		return fmt.Errorf("vector extension unavailable")
	}
	var version int64
	query := `SELECT COALESCE(MAX(version_id),0) FROM (
		SELECT DISTINCT ON (version_id) version_id,is_applied
		FROM goose_db_version ORDER BY version_id,id DESC
	) versions WHERE is_applied`
	if err := db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		return fmt.Errorf("schema version unavailable")
	}
	if version != migrations.LatestVersion {
		return fmt.Errorf("schema version mismatch")
	}
	return nil
}
