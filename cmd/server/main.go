package main

import (
	"database/sql"
	"log"
	"path/filepath"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"merch-app-codex/internal/config"
	"merch-app-codex/internal/report"
	"merch-app-codex/internal/server"
	storage "merch-app-codex/internal/storage/mysql"
)

func main() {
	cfg := config.Load()

	if err := runMigrations(cfg); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	gormDB, err := gorm.Open(gormmysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	repo := storage.NewRepository(gormDB)
	authRepo := storage.NewAuthRepository(gormDB)
	reportService := report.NewService(repo)

	router := server.NewRouter(cfg, repo, authRepo, reportService)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func runMigrations(cfg config.Config) error {
	absPath, err := filepath.Abs(cfg.MigrationsPath)
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+absPath, cfg.MySQLDatabase, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
