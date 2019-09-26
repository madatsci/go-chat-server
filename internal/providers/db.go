package providers

import "github.com/go-pg/pg"

func NewDB(config *Config) (*pg.DB, error) {
	conn := pg.Connect(&pg.Options{
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
		Database: config.DatabaseName,
	})

	if _, err := conn.ExecOne("SELECT 1"); err != nil {
		return nil, err
	}

	return conn, nil
}
