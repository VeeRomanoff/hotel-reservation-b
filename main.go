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

	// handlers initialization
	var (
		app        = fiber.New(config)
		apiv1      = app.Group("/api/v1")
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)

		store = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}

		userHandler  = v1.NewUserHandler(userStore)
		hotelHandler = v1.NewHotelHandler(store)
	)

	// user handlers
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUserById)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandlerDeleteUser)

	// hotel handlers
	apiv1.Post("/hotel", hotelHandler.HandleInsertHotel)
	apiv1.Get("hotel/:id", hotelHandler.HandleGetHotelById)
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRoomsByHotelID)

	app.Listen(*listenAddr)
}
