package providers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"med-chat-bot/internal/dtos"
	"med-chat-bot/internal/errors"
	"med-chat-bot/internal/ginServer"
	"med-chat-bot/internal/meta"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/pkg/db"
	"net/http"
	"strings"
)

func newServerConfig() *ginServer.Config {
	return &ginServer.Config{
		Addr: viper.GetString(cfg.ConfigKeyHttpAddress),
		Port: viper.GetInt64(cfg.ConfigKeyHttpPort),
	}
}

func newErrorParserConfig() *errors.ErrorParserConfig {
	staticErrorCfgPath := "./statics/errors.toml"
	return &errors.ErrorParserConfig{PathConfigError: staticErrorCfgPath}
}

func newGinEngine() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, dtos.Response{
			Meta: meta.Meta{
				Code:    http.StatusNotFound,
				Message: "Page not",
			}})
	})

	return r
}

func newMySQLConnection() *db.DB {
	_db, err := db.Connect(&db.Config{
		Driver:   db.DriverMySQL,
		Username: viper.GetString(cfg.ConfigKeyDBMySQLUsername),
		Password: viper.GetString(cfg.ConfigKeyDBMySQLPassword),
		Host:     viper.GetString(cfg.ConfigKeyDBMySQLHost),
		Port:     viper.GetInt64(cfg.ConfigKeyDBMySQLPort),
		Database: viper.GetString(cfg.ConfigKeyDBMySQLDatabase),
	})
	if err != nil {
		log.Fatalf("Connecting to MySQL DB: %v", err)
	}
	return _db
}

func newMySQLUserTrackingConnection() *db.DB {
	_db, err := db.Connect(&db.Config{
		Driver: db.DriverMySQL,
		//Username: viper.GetString(cfg.ConfigKeyDBMySQLUsername),
		//Password: viper.GetString(cfg.ConfigKeyDBMySQLPassword),
		//Host:     viper.GetString(cfg.ConfigKeyDBMySQLHost),
		//Port:     viper.GetInt64(cfg.ConfigKeyDBMySQLPort),
		//Database: viper.GetString(cfg.ConfigKeyDBMySQLTrackingDatabase),
		Username: "root",
		Password: "tungoccho123",
		Host:     "localhost",
		Port:     3306,
		Database: "tracking_test",
	})
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
