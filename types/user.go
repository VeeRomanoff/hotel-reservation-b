package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost      = 10
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

// UserDTO REQUEST SCOPE
type UserDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u *UserDTO) Validate() map[string]string {
	errors := map[string]string{}
	if len(u.FirstName) < minFirstNameLen {
		errors["first_name"] = fmt.Sprintf("first name length should be at least %d characters", minFirstNameLen)
	}
	if len(u.LastName) < minLastNameLen {
		errors["last_name"] = fmt.Sprintf("last name length should be at least %d characters", minLastNameLen)
	}
	if len(u.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters long", minPasswordLen)
	}
	if !isValidEmail(u.Email) {
		errors["email"] = fmt.Sprintf("email is invalid!!!!!!!")
	}
	return errors
}

func isValidEmail(e string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(e)
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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
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
