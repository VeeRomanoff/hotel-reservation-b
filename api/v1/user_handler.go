package v1

import (
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(ctx *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Bond",
	}
	return ctx.JSON(u)
}

func HandleGetUserById(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]interface{}{
		"success": true,
		"James":   "you were supposed to find the user by id dumass",
	})
}
