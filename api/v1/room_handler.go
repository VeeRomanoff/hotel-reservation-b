package v1

import (
	"context"
	"fmt"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type BookRoomParams struct {
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	NumPersons int       `json:"numPersons"`
	// roomId - ? we don't need since we're passing room into path variables
	// userId - ? each time user authenticates an endpoint we validate the token and fetch their id
}

// validation logic
func (p *BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.StartDate) || now.After(p.EndDate) {
		return fmt.Errorf("can't book a room in the past")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(ctx *fiber.Ctx) error {
	// parsing dto
	var params BookRoomParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	// validating the room dto
	if err := params.validate(); err != nil {
		return err
	}
	// casting roomId into primitive object id
	roomId := ctx.Params("id")
	roomOID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return err
	}
	// fetching current user
	user, ok := ctx.Context().UserValue("user").(*types.User) // this is where we fetch user from since they're authenticated request.
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type:    "error",
			Message: "internal server error", // couldn't fetch user from current session
		})
	}

	ok, err = h.isRoomAvailableForBooking(ctx.Context(), roomOID, params)
	if err != nil {
		return err
	}
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(genericResponse{
			Type:    "error",
			Message: fmt.Sprintf("room %s is already booked", roomId),
		})
	}
	// creating a new booking
	booking := types.Booking{
		RoomID:     roomOID,
		UserID:     user.ID,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		NumPersons: params.NumPersons,
	}
	// inserting a new booking since roomhandler struct has access to store whereas store has access to bookingStore interface
	inserted, err := h.store.Booking.InsertBooking(ctx.Context(), &booking)
	if err != nil {
		return err
	}
	return ctx.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomOID primitive.ObjectID, params BookRoomParams) (bool, error) {
	filter := bson.M{
		// for only the specific room
		"roomID": roomOID,
		"startDate": bson.M{
			"$gte": params.StartDate,
		},
		"endDate": bson.M{
			"$lte": params.EndDate,
		},
	}
	// checking if a room is already booked. the length of bookings in a slice is greater than 0, then it's already booked
	bookings, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
