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

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	// HANDLERS INITIALIZATION
	userHandler := v1.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUserById)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandlerDeleteUser)
	app.Listen(*listenAddr)
}
