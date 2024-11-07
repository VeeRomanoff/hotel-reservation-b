package v1

import (
	"fmt"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
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

type RoomHandler struct {
	roomStore *db.Store
	//userStore *db.UserStore --- ?
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		roomStore: store,
	}
}

func (h *RoomHandler) HandleBookRoom(ctx *fiber.Ctx) error {
	var params BookRoomParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	roomId := ctx.Params("id")
	roomOID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return err
	}
	user, ok := ctx.Context().UserValue("user").(*types.User) // this is where we fetch user from since they're authenticated request.
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type:    "error",
			Message: "internal server error",
		})
	}
	booking := types.Booking{
		RoomID:     roomOID,
		UserID:     user.ID,
		StartDate:  params.StartDate, // TODO: fix time format problem
		EndDate:    params.EndDate,
		NumPersons: params.NumPersons,
	}
	fmt.Printf("%+v\n", booking)
	return nil
}
