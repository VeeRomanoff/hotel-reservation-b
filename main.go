package main

import (
	"context"
	"flag"
	v1 "github.com/VeeRomanoff/hotel-reservation/api/v1"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var listenAddr *string = flag.String("listenAddr", ":4000", "The listen address of the API SERVER")

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{
			"error": err.Error(),
		})
	},
}

func main() {

	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// USER HANDLERS INITIALIZATION
	userHandler := v1.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	// user handlers
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUserById)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandlerDeleteUser)

	// HOTEL HANDLERS INITIALIZATION
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	hotelHandler := v1.NewHotelHandler(db.NewMongoHotelStore(client), roomStore)

	// hotel handlers
	apiv1.Post("/hotel", hotelHandler.HandleInsertHotel)
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRoomsByHotelID)

	app.Listen(*listenAddr)
}
