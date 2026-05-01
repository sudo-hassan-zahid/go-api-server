package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
)

func bindAndValidate[T any](c *fiber.Ctx) (T, bool, error) {
	var req T
	if err := c.BodyParser(&req); err != nil {
		return req, false, err
	}
	if ok := utils.ValidateStruct(c, &req); !ok {
		return req, false, nil
	}
	return req, true, nil
}

func uuidParam(c *fiber.Ctx, name string) (uuid.UUID, error) {
	rawID := strings.TrimSpace(c.Params(name))
	if rawID == "" {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, name+" is required")
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "invalid "+name+" format")
	}

	return id, nil
}
