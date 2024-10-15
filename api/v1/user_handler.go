package v1

import (
	"errors"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlerDeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := h.userStore.DeleteUser(ctx.Context(), id); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"deleted": id})
}

func (h *UserHandler) HandlePutUser(ctx *fiber.Ctx) error {
	var (
		// update bson.M
		values types.UpdateUserDTO
		userID = ctx.Params("id")
	)
	if err := ctx.BodyParser(&values); err != nil {
		return err
	}
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	if err := h.userStore.PutUser(ctx.Context(), filter, values); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"updated": userID})
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

	uCreated, err := h.userStore.InsertUser(ctx.Context(), u)
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.JSON(map[string]string{"error": "not found"})
		}
	}
	return ctx.JSON(user)
}
