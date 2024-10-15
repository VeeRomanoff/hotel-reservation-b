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
	//todo: внедрить рум в отель, но для этого мы будем создавать storage для отеля
	hotel := types.Hotel{
		Name:     "Five seasons",
		Location: "France",
	}

	room := types.Room{
		Type:      types.SingleRoomType,
		BasePrice: 299.9,
	}
	_ = room
	fmt.Println("seeding the database...")
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	room.HotelID = insertedHotel.ID
	fmt.Println("inserted hotel:", insertedHotel)
}
