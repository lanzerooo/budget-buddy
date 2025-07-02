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

	// Проверка и создание таблицы categories
	var tableExists bool
	err = db.QueryRow(`SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'categories'
    )`).Scan(&tableExists)
	if err != nil {
		logger.Error("Failed to check if categories table exists: ", err)
		return err
	}
	if !tableExists {
		_, err = db.Exec(`
            CREATE TABLE categories (
                id SERIAL PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                type VARCHAR(50) NOT NULL CHECK (type IN ('income', 'expense')),
                UNIQUE(name, type)
            )
        `)
		if err != nil {
			logger.Error("Failed to create categories table: ", err)
			return err
		}
		logger.Info("Categories table created successfully")
	}

	// Проверка и создание таблицы subcategories
	err = db.QueryRow(`SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'subcategories'
    )`).Scan(&tableExists)
	if err != nil {
		logger.Error("Failed to check if subcategories table exists: ", err)
		return err
	}
	if !tableExists {
		_, err = db.Exec(`
            CREATE TABLE subcategories (
                id SERIAL PRIMARY KEY,
                category_id INTEGER REFERENCES categories(id),
                name VARCHAR(255) NOT NULL
            )
        `)
		if err != nil {
			logger.Error("Failed to create subcategories table: ", err)
			return err
		}
		logger.Info("Subcategories table created successfully")
	}

	// Проверка и создание таблицы incomes
	err = db.QueryRow(`SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'incomes'
    )`).Scan(&tableExists)
	if err != nil {
		logger.Error("Failed to check if incomes table exists: ", err)
		return err
	}
	if !tableExists {
		_, err = db.Exec(`
            CREATE TABLE incomes (
                id SERIAL PRIMARY KEY,
                user_id INTEGER NOT NULL,
                amount DECIMAL(10,2) NOT NULL,
                category_id INTEGER REFERENCES categories(id),
                subcategory_id INTEGER REFERENCES subcategories(id),
                description TEXT,
                tags TEXT[],
                date DATE NOT NULL,
                note TEXT
            )
        `)
		if err != nil {
			logger.Error("Failed to create incomes table: ", err)
			return err
		}
		logger.Info("Incomes table created successfully")
	} else {
		// Проверка и добавление столбцов для incomes
		var columnExists bool
		err = db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = 'incomes' 
            AND column_name = 'subcategory_id'
        )`).Scan(&columnExists)
		if err != nil {
			logger.Error("Failed to check if subcategory_id column exists in incomes: ", err)
			return err
		}
		if !columnExists {
			_, err = db.Exec(`ALTER TABLE incomes ADD COLUMN subcategory_id INTEGER REFERENCES subcategories(id)`)
			if err != nil {
				logger.Error("Failed to add subcategory_id column to incomes: ", err)
				return err
			}
			logger.Info("Added subcategory_id column to incomes")
		}

		err = db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = 'incomes' 
            AND column_name = 'description'
        )`).Scan(&columnExists)
		if err != nil {
			logger.Error("Failed to check if description column exists in incomes: ", err)
			return err
		}
		if !columnExists {
			_, err = db.Exec(`ALTER TABLE incomes ADD COLUMN description TEXT`)
			if err != nil {
				logger.Error("Failed to add description column to incomes: ", err)
				return err
			}
			logger.Info("Added description column to incomes")
		}

		err = db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = 'incomes' 
            AND column_name = 'tags'
        )`).Scan(&columnExists)
		if err != nil {
			logger.Error("Failed to check if tags column exists in incomes: ", err)
			return err
		}
		if !columnExists {
			_, err = db.Exec(`ALTER TABLE incomes ADD COLUMN tags TEXT[]`)
			if err != nil {
				logger.Error("Failed to add tags column to incomes: ", err)
				return err
			}
			logger.Info("Added tags column to incomes")
		}
	}

	// Проверка и создание таблицы expenses
	err = db.QueryRow(`SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'expenses'
    )`).Scan(&tableExists)
	if err != nil {
		logger.Error("Failed to check if expenses table exists: ", err)
		return err
	}
	if !tableExists {
		_, err = db.Exec(`
            CREATE TABLE expenses (
                id SERIAL PRIMARY KEY,
                user_id INTEGER NOT NULL,
                amount DECIMAL(10,2) NOT NULL,
                category_id INTEGER REFERENCES categories(id),
                subcategory_id INTEGER REFERENCES subcategories(id),
                description TEXT,
                tags TEXT[],
                date DATE NOT NULL,
                note TEXT
            )
        `)
		if err != nil {
			logger.Error("Failed to create expenses table: ", err)
			return err
		}
		logger.Info("Expenses table created successfully")
	} else {
		// Проверка и добавление столбцов для expenses
		var columnExists bool
		err = db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = 'expenses' 
            AND column_name = 'subcategory_id'
        )`).Scan(&columnExists)
		if err != nil {
			logger.Error("Failed to check if subcategory_id column exists in expenses: ", err)
			return err
		}
		if !columnExists {
			_, err = db.Exec(`ALTER TABLE expenses ADD COLUMN subcategory_id INTEGER REFERENCES subcategories(id)`)
			if err != nil {
				logger.Error("Failed to add subcategory_id column to expenses: ", err)
				return err
			}
			logger.Info("Added subcategory_id column to expenses")
		}

		err = db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = 'expenses' 
            AND column_name = 'description'
        )`).Scan(&columnExists)
		if err != nil {
			logger.Error("Failed to check if description column exists in expenses: ", err)
			return err
		}
		if !columnExists {
			_, err = db.Exec(`ALTER TABLE expenses ADD COLUMN description TEXT`)
			if err != nil {
				logger.Error("Failed to add description column to expenses: ", err)
				return err
			}
			logger.Info("Added description column to expenses")
		}

		err = db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.columns 
            WHERE table_schema = 'public' 
            AND table_name = 'expenses' 
            AND column_name = 'tags'
        )`).Scan(&columnExists)
		if err != nil {
			logger.Error("Failed to check if tags column exists in expenses: ", err)
			return err
		}
		if !columnExists {
			_, err = db.Exec(`ALTER TABLE expenses ADD COLUMN tags TEXT[]`)
			if err != nil {
				logger.Error("Failed to add tags column to expenses: ", err)
				return err
			}
			logger.Info("Added tags column to expenses")
		}
	}

	// Проверка и создание таблицы goals
	err = db.QueryRow(`SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'goals'
    )`).Scan(&tableExists)
	if err != nil {
		logger.Error("Failed to check if goals table exists: ", err)
		return err
	}
	if !tableExists {
		_, err = db.Exec(`
            CREATE TABLE goals (
                id SERIAL PRIMARY KEY,
                user_id INTEGER NOT NULL,
                name VARCHAR(255) NOT NULL,
                target_amount DECIMAL(10,2) NOT NULL,
                current_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
                deadline DATE NOT NULL,
                created_at TIMESTAMP NOT NULL
            )
        `)
		if err != nil {
			logger.Error("Failed to create goals table: ", err)
			return err
		}
		logger.Info("Goals table created successfully")
	}

	logger.Info("Finance migrations executed successfully")
	return nil
}
