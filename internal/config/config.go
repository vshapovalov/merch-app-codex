package config

import (
	"fmt"
	"os"
)

// Config contains application level configuration values loaded from environment variables.
type Config struct {
	Port           string
	MySQLHost      string
	MySQLPort      string
	MySQLUser      string
	MySQLPassword  string
	MySQLDatabase  string
	MigrationsPath string
}

// Load reads configuration from environment variables and applies sensible defaults.
func Load() Config {
	cfg := Config{
		Port:           getEnv("PORT", "8080"),
		MySQLHost:      getEnv("MYSQL_HOST", "127.0.0.1"),
		MySQLPort:      getEnv("MYSQL_PORT", "3306"),
		MySQLUser:      getEnv("MYSQL_USER", "root"),
		MySQLPassword:  os.Getenv("MYSQL_PASSWORD"),
		MySQLDatabase:  getEnv("MYSQL_DATABASE", "merch"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "db/migrations"),
	}

	return cfg
}

// DSN returns a Data Source Name that can be used by database/sql and GORM.
func (c Config) DSN() string {
	credentials := c.MySQLUser
	if c.MySQLPassword != "" {
		credentials = fmt.Sprintf("%s:%s", c.MySQLUser, c.MySQLPassword)
	}

	return fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&multiStatements=true",
		credentials,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDatabase,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
