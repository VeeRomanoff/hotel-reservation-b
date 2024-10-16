package main

import (
	"context"
	"fmt"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client)

	hotel := types.Hotel{
		Name:     "Five seasons",
		Location: "France",
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 299.9,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 499.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 1999.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 799.9,
		},
	}
	fmt.Println("seeding the database...")
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted hotel:", insertedHotel)

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("inserted room: ", insertedRoom)
	}
}
