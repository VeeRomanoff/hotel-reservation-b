package main

import (
	"context"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	// RESETTING DB EACH TIME APP STARTS
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromDTO(types.UserDTO{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  "supersecurepassword",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			Size:      "small",
			BasePrice: 299.9,
		},
		{
			Type:      types.DoubleRoomType,
			Size:      "normal",
			BasePrice: 499.9,
		},
		{
			Type:      types.DeluxeRoomType,
			Size:      "kingsize",
			BasePrice: 1999.9,
		},
		{

			Type:      types.SeaSideRoomType,
			Size:      "normal",
			BasePrice: 799.9,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// SEED SOME HOTELS IN
	seedHotel("Bellucia", "France", 3)
	seedHotel("Five seasons", "The Netherlands", 4)
	seedHotel("Sultan Ahmet", "Turkey", 1)

	// SEED SOME USERS IN
	seedUser("James", "In the bottom", "james@inthebottom.com")
}
