package providers

import (
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

func NewDB(config *viper.Viper) (*pg.DB, error) {
	config.SetDefault("db.user", "postgres")
	config.SetDefault("db.password", "password")
	config.SetDefault("db.db", "go_chat_server")
	config.SetDefault("db.addr", "addr")

	conn := pg.Connect(&pg.Options{
		User:     config.GetString("db.user"),
		Password: config.GetString("db.password"),
		Database: config.GetString("db.db"),
		Addr:     config.GetString("db.addr"),
	})

	if _, err := conn.ExecOne("SELECT 1"); err != nil {
		return nil, err
	}

	return conn, nil
}
