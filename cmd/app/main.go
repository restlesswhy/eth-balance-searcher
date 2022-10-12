package main

import (
	"context"
	"fmt"
	"log"

	"github.com/restlesswhy/eth-balance-searcher/config"
	"github.com/restlesswhy/eth-balance-searcher/internal/server"
	"github.com/restlesswhy/eth-balance-searcher/pkg/logger"
	"github.com/restlesswhy/eth-balance-searcher/pkg/redis"
)

// @title ETH balance searcher
// @version 2.0
// @description Service

// @contact.name German Generalov
// @contact.url http://github.com/restlesswhy

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /api/v1/
// @schemes http
func main() {
	log.Println("Starting microservice")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	redis, err := redis.New(cfg, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.Named(fmt.Sprintf(`(%s)`, cfg.ServiceName))
	appLogger.Infof("CFG: %+v", cfg)

	if err := server.New(appLogger, cfg, redis).Run(); err != nil {
		appLogger.Fatal(err)
	}
}
