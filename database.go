package main

import (
	"database/sql"
	"fmt"
	"github.com/blang/semver/v4"
	_ "github.com/lib/pq"
)

const databaseDriver = "postgres"

type DatabaseConfig struct {
	Host       string
	Port       int
	Database   string
	SearchPath string
	User       string
	Password   string
	Version    string
}

func newDatabase(cfg DatabaseConfig) (*sql.DB, error) {
	fail := func(err error) (*sql.DB, error) {
		return nil, fmt.Errorf("newDatabase: %w", err)
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%d dbname=%s search_path=%s user=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SearchPath,
		cfg.User,
		cfg.Password,
	)

	db, err := sql.Open(databaseDriver, dataSource)
	if err != nil {
		return fail(fmt.Errorf("failed to connect to database: %w", err))
	}

	if err := db.Ping(); err != nil {
		return fail(fmt.Errorf("failed to pingHandler database: %w", err))
	}

	configVersion, err := semver.Make(cfg.Version)
	if err != nil {
		return fail(fmt.Errorf("config has invalid version number: %w", err))
	}
	databaseVersion, err := getMigrationVersion(db)
	if err != nil {
		return fail(fmt.Errorf("failed to determine database migration version: %w", err))
	}
	if !configVersion.Equals(*databaseVersion) {
		return fail(fmt.Errorf("require database migration: got version %v, require version %v", databaseVersion, configVersion))
	}

	return db, nil
}

func getMigrationVersion(db *sql.DB) (*semver.Version, error) {
	fail := func(err error) (*semver.Version, error) {
		return nil, fmt.Errorf("getMigrationVersion: %w", err)
	}
	query := `SELECT version FROM flyway_schema_history WHERE success IS TRUE ORDER BY installed_rank DESC LIMIT 1`
	var versionStr string
	if err := db.QueryRow(query).Scan(&versionStr); err != nil {
		return fail(err)
	}
	if version, err := semver.Make(versionStr); err != nil {
		return fail(fmt.Errorf("database has invalid version number: %w", err))
	} else {
		return &version, nil
	}
}
