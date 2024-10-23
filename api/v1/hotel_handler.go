package v1

import (
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(map[string]any{"hotels": hotels, "totalCount": len(hotels)})
}

func (h *HotelHandler) HandleGetRoomsByHotelID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	rooms, err := h.store.Room.GetRooms(ctx.Context(), bson.M{"hotelID": oid})
	if err != nil {
		return err
	}
	return ctx.JSON(map[string]any{"rooms": rooms, "totalCount": len(rooms)})
}

func (h *HotelHandler) HandleInsertHotel(ctx *fiber.Ctx) error {
	var hotelDTO types.HotelDTO
	if err := ctx.BodyParser(&hotelDTO); err != nil {
		return err
	}
	hotel, err := types.NewHotelFromDTO(hotelDTO)
	if err != nil {
		return err
	}
	insertedHotel, err := h.store.Hotel.InsertHotel(ctx.Context(), hotel)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedHotel)
}

func (h *HotelHandler) HandleGetHotelById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	hotel, err := h.store.Hotel.GetHotelById(ctx.Context(), oid)
	if err != nil {
		return err
	}
	return ctx.JSON(hotel)
}
