package v1

import (
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	// DEPENDENCY
	roomStore db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(ctx.Context(), nil)
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
	rooms, err := h.roomStore.GetRooms(ctx.Context(), bson.M{"hotelID": oid})
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
	insertedHotel, err := h.hotelStore.InsertHotel(ctx.Context(), hotel)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedHotel)
}
