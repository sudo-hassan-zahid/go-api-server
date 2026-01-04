package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateStruct(c *fiber.Ctx, s interface{}) bool {
	err := validate.Struct(s)
	if err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Error()
		}
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
		return false
	}
	return true
}
