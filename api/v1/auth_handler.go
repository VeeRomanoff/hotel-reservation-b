package v1

import (
	"errors"
	"fmt"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/VeeRomanoff/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

// HandleAuthenticate A handler should only do:
//   - a serialization of the incoming req
//   - do some data fetching from db
//   - call some business logic!!
//   - return the data back to user (@Controller analogue in java)
func (h *AuthHandler) HandleAuthenticate(ctx *fiber.Ctx) error {
	var params AuthParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("Invalid credentials")
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("invalid email or password")
	}
	token := createTokenFromUser(user)

	resp := AuthResponse{
		User:  user,
		Token: token,
	}
	return ctx.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	expires := time.Now().Add(time.Hour * 3).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret")
	}
	return tokenStr
}
