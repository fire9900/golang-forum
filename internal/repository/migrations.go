package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDB выполняет миграции базы данных
func MigrateDB(db *sql.DB, migrationsPath string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("не удалось создать экземпляр драйвера mysql: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("не удалось создать экземпляр migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("не удалось выполнить миграции: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("не удалось получить версию миграции: %w", err)
	}

	log.Printf("Текущая версия миграции: %d, dirty: %v\n", version, dirty)
	return nil
}
