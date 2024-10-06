package main

import (
	"flag"
	v1 "github.com/VeeRomanoff/hotel-reservation/api/v1"
	"github.com/gofiber/fiber/v2"
)

var listenAddr *string = flag.String("listenAddr", ":4000", "The listen address of the API SERVER")

func main() {
	app := fiber.New()
	flag.Parse()
	apiv1 := app.Group("api/v1/")

	apiv1.Get("/user", v1.HandleGetUsers)
	apiv1.Get("/user/:id", v1.HandleGetUserById)
	app.Listen(*listenAddr)

}
