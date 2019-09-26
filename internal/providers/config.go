package providers

import (
	"os"
)

type Config struct {
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	ListenAddr       string
}

// NewConfig creates new application config
func NewConfig() (*Config, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	addr := os.Getenv("ADDR")

	if len(dbUser) == 0 {
		dbUser = "postgres"
	}

	if len(dbPassword) == 0 {
		dbPassword = "password"
	}

	if len(dbName) == 0 {
		dbName = "go_chat_server"
	}

	if len(addr) == 0 {
		addr = ":9001"
	}

	return &Config{
		DatabaseUser:     dbUser,
		DatabasePassword: dbPassword,
		DatabaseName:     dbName,
		ListenAddr:       addr,
	}, nil
}
