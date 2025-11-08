package postgres

import (
	"database/sql"
	"fmt"
	"github.com/delapaska/EgorExplain/pkg/config"
	"github.com/delapaska/EgorExplain/pkg/postgres/migration"
	"log"
)

func NewPostgresSQLStorage() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Envs.Host, config.Envs.DBPort,
		config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	migrator := migration.NewMigrator(db, migrationsDir)
	return migrator.Up()
}
