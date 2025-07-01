package migrations

import (
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"
	"database/sql"
	"embed"

	_ "github.com/lib/pq"
)

//go:embed *.sql
var migrations embed.FS

func RunMigrations(cfg *config.Config) error {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
		return err
	}
	defer db.Close()

	// Чтение и выполнение миграции
	migrationSQL, err := migrations.ReadFile("001_create_users_table.sql")
	if err != nil {
		logger.Error("Failed to read migration file: ", err)
		return err
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		logger.Error("Failed to execute migration: ", err)
		return err
	}

	logger.Info("Migrations executed successfully")
	return nil
}
