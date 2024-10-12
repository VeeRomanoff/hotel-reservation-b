package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// TODO собирается учить как отеделять доменые сущности от реквестовых сущностей.

// UserDTO REQUEST SCOPE
type UserDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// User DOMAIN SCOPE
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"first_name" json:"first_name"`
	LastName          string             `bson:"last_name" json:"last_name"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encrypted_password" json:"-"` // SINCE WE DON`T WANT TO SHARE THIS WHEN FETCHING USER
}

func NewUserFromDTO(u UserDTO) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Email:             u.Email,
		EncryptedPassword: string(passwordHash),
	}, nil
}
