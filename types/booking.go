package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"roomID,omitempty" json:"userID,omitempty"`
	RoomID     primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"` // already has got the hotelid. no need to embed hotelid to booking struct
	NumPersons int                `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	StartDate  time.Time          `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate    time.Time          `bson:"endDate,omitempty" json:"endDate,omitempty"`
}
