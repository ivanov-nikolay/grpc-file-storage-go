package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type MigrationManager struct {
	db *sql.DB
}

func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{
		db: db,
	}
}

func (m *MigrationManager) CreateMigrationsTable() error {
	query := `CREATE TABLE IF NOT EXISTS schema_migrations (
    				version VARCHAR(255) PRIMARY KEY,
    				applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)
			`
	_, err := m.db.Exec(query)

	return err
}

func (m *MigrationManager) IsMigrationApplied(version string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM schema_migrations WHERE version = $1`
	err := m.db.QueryRow(query, version).Scan(&count)

	return count > 0, err
}

func (m *MigrationManager) MarkMigrationApplied(version string) error {
	query := `INSERT INTO schema_migrations (version) VALUES ($1)`
	_, err := m.db.Exec(query, version)

	return err
}

func (m *MigrationManager) RunMigrations(migrationsDir string) error {
	if err := m.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			version := file.Name()

			applied, err := m.IsMigrationApplied(version)
			if err != nil {
				return fmt.Errorf("failed to check migration status: %v", err)
			}
			if applied {
				log.Printf("migration %s is already applied", version)

				continue
			}

			migrationPath := filepath.Join(migrationsDir, version)
			content, err := os.ReadFile(migrationPath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %v", version, err)
			}

			tx, err := m.db.Begin()
			if err != nil {
				return fmt.Errorf("failed to start transaction: %v", err)
			}

			if _, err := tx.Exec(string(content)); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to execute migration: %v", err)
			}

			if err := m.MarkMigrationApplied(version); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to mark migration as applied: %v", err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit transaction: %v", err)
			}
			log.Printf("migration %s applied successfully", version)
		}
	}
	log.Println("all migration applied successfully")

	return nil
}
