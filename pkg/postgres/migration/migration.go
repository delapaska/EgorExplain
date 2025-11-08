package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) getAppliedMigrations() (map[int]bool, error) {
	query := "SELECT version FROM schema_migrations"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func (m *Migrator) loadMigrations() ([]Migration, error) {
	var migrations []Migration

	err := filepath.Walk(m.migrationsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".up.sql") {
			return nil
		}

		fileName := info.Name()
		parts := strings.Split(fileName, "_")
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration file name: %s", fileName)
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid version in migration file: %s", fileName)
		}

		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".up.sql")

		upContent, err := os.ReadFile(filepath.Join(m.migrationsDir, fileName))
		if err != nil {
			return fmt.Errorf("failed to read up migration file: %w", err)
		}

		downFileName := fmt.Sprintf("%03d_%s.down.sql", version, name)
		downContent, err := os.ReadFile(filepath.Join(m.migrationsDir, downFileName))
		if err != nil {
			return fmt.Errorf("failed to read down migration file: %w", err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(upContent),
			DownSQL: string(downContent),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (m *Migrator) Up() error {
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		if _, err := tx.Exec(migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}

		insertQuery := "INSERT INTO schema_migrations (version) VALUES ($1)"
		if _, err := tx.Exec(insertQuery, migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		fmt.Printf("Applied migration %d successfully\n", migration.Version)
	}

	return nil
}

func (m *Migrator) Down(steps int) error {
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	var toRollback []Migration
	for i := len(migrations) - 1; i >= 0 && len(toRollback) < steps; i-- {
		if applied[migrations[i].Version] {
			toRollback = append(toRollback, migrations[i])
		}
	}

	for _, migration := range toRollback {
		fmt.Printf("Rolling back migration %d: %s\n", migration.Version, migration.Name)

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		if _, err := tx.Exec(migration.DownSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
		}

		deleteQuery := "DELETE FROM schema_migrations WHERE version = $1"
		if _, err := tx.Exec(deleteQuery, migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit rollback %d: %w", migration.Version, err)
		}

		fmt.Printf("Rolled back migration %d successfully\n", migration.Version)
	}

	return nil
}
