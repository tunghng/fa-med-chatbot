package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"med-chat-bot/cmd/medi-query-bot/internal/providers"
	"med-chat-bot/internal/ginServer"
	"med-chat-bot/pkg/cfg"
	"os"
)

func main() {
	providers.BuildContainer()

	if os.Getenv("ENVIRONMENT") == cfg.EnvironmentLocal {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("Preparing and running main application . . . !")

	if err := run(); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}
}

func run() error {
	c := providers.GetContainer()
	if c == nil {
		log.Fatalf("Container hasn't been initialized yet")
	}
	var s ginServer.Server
	if err := c.Invoke(func(_s ginServer.Server) {
		s = _s
	}); err != nil {
		return err
	}

	if err := s.Open(); err != nil {
		return err
	}

	return nil
}
