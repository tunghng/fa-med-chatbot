package providers

import (
	"github.com/spf13/viper"
	"log"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/pkg/db"
	"strings"
)

func newMySQLConnection() *db.DB {
	_db, err := db.Connect(&db.Config{
		Driver: db.DriverMySQL,
		//Username: viper.GetString(cfg.ConfigKeyDBMySQLUsername),
		//Password: viper.GetString(cfg.ConfigKeyDBMySQLPassword),
		//Host:     viper.GetString(cfg.ConfigKeyDBMySQLHost),
		//Port:     viper.GetInt64(cfg.ConfigKeyDBMySQLPort),
		Username: "videdent_tele",
		Password: "Muaxuan2024",
		Host:     "173.252.167.20",
		Port:     3306,
		Database: viper.GetString(cfg.ConfigKeyDBMySQLDatabase)})
	if err != nil {
		log.Fatalf("Connecting to MySQL DB: %v", err)
	}
	return _db
}

func newCfgReader() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}
