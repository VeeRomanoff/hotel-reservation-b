package db

const (
	DBNAME string = "hotel-reservation"
	DBURI  string = "mongodb://localhost:27017"
)

// user store
const userCollection = "users"

// hotel store
const hotelCollection = "hotels"

// room store
const roomCollection = "rooms"

// Store is going to contain all storages inside of it itself
type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
