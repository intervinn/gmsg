package middleware

import (
	"log"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/intervinn/gmsg/api/model"
	"github.com/intervinn/gmsg/api/service"
	"github.com/labstack/echo/v5"
	"github.com/uptrace/bun"
)

func RequireAuth(auth *service.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			tokstr := strings.TrimPrefix(header, "Bearer ")
			tok, err := auth.Verify(tokstr)
			if err != nil {
				return echo.ErrUnauthorized.Wrap(err)
			}

			c.Set("token", tok)
			return next(c)
		}
	}
}

// use before RequireAuth()
func RequireUser(db *bun.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token, ok := c.Get("token").(*jwt.Token)
			if !ok {
				log.Println("token not found")
				return echo.ErrUnauthorized
			}

			sub, err := token.Claims.GetSubject()
			if err != nil {
				log.Println("claims fetch failed:", err)
				return echo.ErrUnauthorized
			}

			id, err := strconv.ParseInt(sub, 10, 64)
			if err != nil {
				log.Println("parsing failed:", err)
				return echo.ErrUnauthorized
			}

			ctx := c.Request().Context()

			user := new(model.User)
			err = db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
			if err != nil {
				log.Println("bun: query failed:", err)
				return echo.ErrUnauthorized
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
