package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/bwmarrin/snowflake"
	"github.com/intervinn/gmsg/api/dto"
	"github.com/intervinn/gmsg/api/middleware"
	"github.com/intervinn/gmsg/api/model"
	"github.com/intervinn/gmsg/api/service"
	"github.com/labstack/echo/v5"
	"github.com/uptrace/bun"
)

type AuthController struct {
	token *service.TokenService
	db    *bun.DB
	node  *snowflake.Node
}

func NewAuthController(e *echo.Echo, token *service.TokenService, db *bun.DB, node *snowflake.Node) *AuthController {
	ac := &AuthController{
		token: token,
		db:    db,
		node:  node,
	}

	e.POST("/auth/register", ac.Register)
	e.POST("/auth/login", ac.Login)
	e.GET("/me", ac.Me, middleware.RequireAuth(token), middleware.RequireUser(db))

	return ac
}

func (ac *AuthController) Register(c *echo.Context) error {
	body := new(dto.RegisterBody)
	if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
		return c.Blob(http.StatusBadRequest, "application/json", resInvalidData)
	}

	ctx := c.Request().Context()

	exists, err := ac.db.NewSelect().Model((*model.User)(nil)).Where("username = ?", body.Username).Exists(ctx)
	if err != nil {
		log.Println("bun: query failed:", err)
		return c.Blob(http.StatusBadRequest, "application/json", resServerError)
	}

	if exists {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Ok:    false,
			Error: "username taken",
		})
	}

	hash, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		log.Println("argon: hashing failed", err)
		return c.Blob(http.StatusBadRequest, "application/json", resServerError)
	}

	user := model.User{
		ID:           ac.node.Generate().Int64(),
		Username:     body.Username,
		PasswordHash: hash,
	}

	_, err = ac.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		log.Println("bun: query failed:", err)
		return c.Blob(http.StatusBadRequest, "application/json", resServerError)
	}

	tok, err := ac.token.Generate(user.ID)
	if err != nil {
		log.Println("token: generation failed:", err)
		return c.Blob(http.StatusBadRequest, "application/json", resServerError)
	}

	return c.JSON(http.StatusOK, dto.Response{
		Ok: true,
		Data: dto.RegisterResult{
			ID:            user.ID,
			Username:      user.Username,
			Authorization: tok,
		},
	})
}

func (ac *AuthController) Login(c *echo.Context) error {
	body := new(dto.LoginBody)
	if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
		return c.Blob(http.StatusBadRequest, "application/json", resInvalidData)
	}

	ctx := c.Request().Context()

	user := new(model.User)
	err := ac.db.NewSelect().Model(user).Where("username = ?", body.Username).Scan(ctx)
	if err != nil {
		log.Println("bun: query failed:", err)
		return c.Blob(http.StatusBadRequest, "application/json", resInvalidData)
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, user.PasswordHash)
	if err != nil {
		log.Println("argon: compare failed:", err)
		return c.Blob(http.StatusBadRequest, "application/json", resServerError)
	}

	if !match {
		return c.Blob(http.StatusBadRequest, "application/json", resInvalidData)
	}

	tok, err := ac.token.Generate(user.ID)
	if err != nil {
		log.Println("token: generation failed:", err)
		return c.Blob(http.StatusBadRequest, "application/json", resServerError)
	}

	return c.JSON(http.StatusOK, dto.Response{
		Ok: true,
		Data: dto.LoginResult{
			ID:            user.ID,
			Username:      user.Username,
			Authorization: tok,
		},
	})
}

func (ac *AuthController) Me(c *echo.Context) error {
	user := c.Get("user").(*model.User)
	return c.JSON(http.StatusOK, dto.Response{
		Ok: true,
		Data: dto.PublicUser{
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			ID:        user.ID,
		},
	})
}
