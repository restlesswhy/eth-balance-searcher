package v1

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/restlesswhy/eth-balance-searcher/pkg/logger"
)

type Service interface {
	GetAddress(ctx context.Context) (string, error)
}

type handler struct {
	log     logger.Logger
	service Service
}

func New(log logger.Logger, service Service) *handler {
	return &handler{log: log, service: service}
}

func (h *handler) getAddress(c *fiber.Ctx) error {
	currency, err := h.service.GetAddress(c.Context())
	if err != nil {
		h.log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "get currency error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   currency,
	})
}
