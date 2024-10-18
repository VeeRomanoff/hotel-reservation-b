package v1

import (
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
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
