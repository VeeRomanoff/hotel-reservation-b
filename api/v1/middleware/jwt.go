package middeware

import (
	"fmt"
	"github.com/VeeRomanoff/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

// JWTAuthenticationDecorator uses decorator pattern since we want to add database operations over this function
func JWTAuthenticationDecorator(userStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fmt.Println("-- jwt auth")
		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok || len(token) == 0 {
			fmt.Println("Token is not in the header")
			return fmt.Errorf("unauthorized")
		}
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		// check token expiration
		if time.Now().Unix() > expires {
			return fmt.Errorf("token is expired")
		}
		userId := claims["id"].(string)
		user, err := userStore.GetUserById(ctx.Context(), userId) // Here we fetch user id from any request in app
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		// Set the current authenticated user to the context (and use it how?)
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}

// each time user requests an authenticated endpoint, we validate the token
func validateToken(tokenstr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized. lol")
		}
		// that's how you get env variable inside go
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, fmt.Errorf("unauthorized. lol")
	}
	if !token.Valid {
		fmt.Println("invalid JWT token")
		return nil, fmt.Errorf("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
