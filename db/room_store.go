package db

import (
	"context"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(roomCollection),
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	// WE NEED TO UPDATE THE HOTEL WITH THIS ROOM ID
	//filter := bson.M{"_id": room.HotelID}
	//update := bson.M{"$set": bson.M{"rooms": room.ID}}
	//
	return room, nil
}
