package providers

import (
	"os"
)

// Config provides all the main parameters for running application
type Config struct {
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	ListenAddr       string
	HashComplexity   int
}

var (
	defaultHashComplexity = 10
	defaultDbUser         = "postgres"
	defaultDbPassword     = "password"
	defaultDbName         = "go_chat_server"
	defaultAddr           = ":9001"
)

// NewConfig creates new application config
func NewConfig() (*Config, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	addr := os.Getenv("ADDR")

	if len(dbUser) == 0 {
		dbUser = defaultDbUser
	}

	if len(dbPassword) == 0 {
		dbPassword = defaultDbPassword
	}

	if len(dbName) == 0 {
		dbName = defaultDbName
	}

	if len(addr) == 0 {
		addr = defaultAddr
	}

	return &Config{
		DatabaseUser:     dbUser,
		DatabasePassword: dbPassword,
		DatabaseName:     dbName,
		ListenAddr:       addr,
		HashComplexity:   defaultHashComplexity,
	}, nil
}
