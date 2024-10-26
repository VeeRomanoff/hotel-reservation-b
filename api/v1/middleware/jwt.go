package middeware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func JWTAuthentication(ctx *fiber.Ctx) error {
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
	return ctx.Next()
}

// todo я не понимаю что происходит. остановились на том что получали токен с claims, потом решили проверять не истек ли токен чекая ключ "expires" у claims. че то не рабоатет ничего

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
