package db

import (
	"context"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	// DEPENDENCY. USE ANYTHING THAT IMPLEMENTS HOTELSTORE
	HotelStore
}

// NewMongoRoomStore второй параметр -- это инжектирование нашей зависимости (DEPENDENCY INJECTION)
func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(roomCollection),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	// ALTERING HOTEL STUFF
	filter := bson.M{"_id": room.HotelID}               // find hotel by room.HotelID
	update := bson.M{"$push": bson.M{"rooms": room.ID}} // push room.ID into []rooms field
	if err := s.HotelStore.PutHotel(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}
