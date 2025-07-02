package migrations

import (
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"
	"database/sql"

	_ "github.com/lib/pq"
)

func RunMigrations(cfg *config.Config) error {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
		return err
	}
	defer db.Close()

	// Проверка существования таблицы users
	var tableExists bool
	err = db.QueryRow(`SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'users'
    )`).Scan(&tableExists)
	if err != nil {
		logger.Error("Failed to check if users table exists: ", err)
		return err
	}

	if !tableExists {
		// Создание таблицы users
		_, err = db.Exec(`
            CREATE TABLE users (
                id SERIAL PRIMARY KEY,
                email VARCHAR(255) UNIQUE NOT NULL,
                password VARCHAR(255) NOT NULL,
                name VARCHAR(255) NOT NULL,
                created_at TIMESTAMP NOT NULL
            )
        `)
		if err != nil {
			logger.Error("Failed to create users table: ", err)
			return err
		}
		logger.Info("Users table created successfully")
	} else {
		logger.Info("Users table already exists")
	}

	logger.Info("User migrations executed successfully")
	return nil
}
