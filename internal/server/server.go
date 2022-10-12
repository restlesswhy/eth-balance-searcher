package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/restlesswhy/eth-balance-searcher/config"
	v1 "github.com/restlesswhy/eth-balance-searcher/internal/delivery/http/v1"
	"github.com/restlesswhy/eth-balance-searcher/internal/integration/getblock"
	"github.com/restlesswhy/eth-balance-searcher/internal/service"
	"github.com/restlesswhy/eth-balance-searcher/pkg/logger"

	"github.com/gofiber/fiber/v2"

	"net/http"
	_ "net/http/pprof"
)

type server struct {
	log   logger.Logger
	cfg   *config.Config
	fiber *fiber.App
}

func New(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, fiber: fiber.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	getBlockRPC := getblock.New(s.cfg)

	service := service.New(s.log, getBlockRPC)
	controller := v1.New(s.log, service)
	controller.SetupRoutes(s.fiber)

	go func() {
		if err := s.runHttp(); err != nil {
			s.log.Errorf("(runHttp) err: %v", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: %v", s.getMicroserviceName(), s.cfg.Http.Port)

	go func() {
		s.log.Error(http.ListenAndServe(":6060", nil))
	}()

	<-ctx.Done()

	if err := s.fiber.Shutdown(); err != nil {
		s.log.Warnf("(Shutdown) err: %v", err)
		return err
	}

	s.log.Info("Service gracefully closed.")
	return nil
}

func (s server) getMicroserviceName() string {
	return fmt.Sprintf("(%s)", strings.ToUpper(s.cfg.ServiceName))
}
