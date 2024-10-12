package v1

import (
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

func (h *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var userDto types.UserDTO
	if err := ctx.BodyParser(&userDto); err != nil {
		return err
	}
	if errors := userDto.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}
	u, err := types.NewUserFromDTO(userDto)
	if err != nil {
		return err
	}

	uCreated, err := h.userStore.PostUser(ctx.Context(), u)
	if err != nil {
		return err
	}
	return ctx.JSON(uCreated)
}

func (h *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(users)
}

func (h *UserHandler) HandleGetUserById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := h.userStore.GetUserById(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(user)
}
