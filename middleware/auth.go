package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sebastian-nunez/golang-search-engine/config"
	"github.com/sebastian-nunez/golang-search-engine/utils"
)

type AdminClaims struct {
	ID                   string `json:"id"`
	User                 string `json:"user"`
	jwt.RegisteredClaims `json:"claims"`
}

func WithAuth(c *fiber.Ctx) error {
	signedToken := c.Cookies(utils.AdminCookie)
	if signedToken == "" {
		return c.Redirect("/login", fiber.StatusFound)
	}

	token, err := jwt.ParseWithClaims(signedToken, &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.SecretKey), nil
	})
	if err != nil {
		return c.Redirect("/login", fiber.StatusFound)
	}

	if _, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return c.Next()
	}

	return c.Redirect("/login", fiber.StatusFound)
}
