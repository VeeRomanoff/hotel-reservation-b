package v1

import (
	"context"
	"fmt"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Bond",
	}
	return ctx.JSON(u)
}

func (h *UserHandler) HandleGetUserById(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
		c  = context.Background()
	)
	user, err := h.userStore.GetUserById(c, id)
	if err != nil {
		return err
	}
	fmt.Println(user)
	return ctx.JSON(user)
}
