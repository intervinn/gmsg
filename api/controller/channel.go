package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/intervinn/gmsg/api/dto"
	"github.com/intervinn/gmsg/api/middleware"
	"github.com/intervinn/gmsg/api/model"
	"github.com/intervinn/gmsg/api/service"
	"github.com/labstack/echo/v5"
	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
)

var (
	resInvalidData = []byte("{\"ok\": false, \"error\": \"invalid data\"}")
	resEmptyOK     = []byte("{\"ok\": true}")
	resServerError = []byte("{\"ok\": false, \"error\": \"internal server error\"}")
)

type ChannelController struct {
	node *snowflake.Node
	nc   *nats.Conn
	db   *bun.DB
}

func NewChannelController(e *echo.Echo, nc *nats.Conn, node *snowflake.Node, db *bun.DB, token *service.TokenService) *ChannelController {
	cc := &ChannelController{}
	cc.nc = nc
	cc.node = node
	cc.db = db

	e.POST("/channels/:channel/messages", cc.SendMessage, middleware.RequireAuth(token), middleware.RequireUser(db))
	e.POST("/guilds/:guild/channels", cc.CreateChannel, middleware.RequireAuth(token))
	e.POST("/guilds", cc.CreateGuild, middleware.RequireAuth(token), middleware.RequireUser(db))

	return cc
}

func (cc *ChannelController) SendMessage(c *echo.Context) error {
	user := c.Get("user").(*model.User)
	body := new(dto.CreateMessageBody)

	ch := c.Param("channel")
	chanID, err := strconv.ParseInt(ch, 10, 64)
	if ch == "" || err != nil {
		return c.Blob(http.StatusBadRequest, "application/json", resInvalidData)
	}

	if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	channel := new(model.Channel)
	err = cc.db.NewSelect().
		Model(channel).
		Where("id = ?", chanID).
		Column("guild_id").
		Scan(ctx)

	if err != nil || channel.GuildID == nil {
		log.Println("bun: query failed:", err)
		return c.Blob(http.StatusInternalServerError, "application/json", resServerError)
	}

	msg := model.Message{
		ID:        cc.node.Generate().Int64(),
		Content:   body.Content,
		AuthorID:  user.ID,
		GuildID:   channel.GuildID,
		ChannelID: chanID,
	}

	_, err = cc.db.NewInsert().Model(&msg).Exec(ctx)
	if err != nil {
		log.Println("bun: query failed:", err)
		return c.Blob(http.StatusInternalServerError, "application/json", resServerError)
	}

	b, err := json.Marshal(&dto.EventMessage{
		Type: "message_create",
		Data: &msg,
	})
	if err != nil {
		return err
	}

	err = cc.nc.Publish(fmt.Sprintf("guild.%v.message", *channel.GuildID), b)
	if err != nil {
		log.Println("failed to publish to nats:", err)
	}

	return c.JSON(http.StatusOK, &msg)
}

func (cc *ChannelController) CreateChannel(c *echo.Context) error {
	guild := c.Param("guild")
	if guild == "" {
		return c.Blob(http.StatusBadRequest, "application/json", resInvalidData)
	}

	guildID, err := strconv.ParseInt(guild, 10, 64)
	if err != nil {
		return c.Blob(http.StatusBadRequest, "application/json", resInvalidData)
	}

	body := new(dto.CreateChannelBody)
	if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
		return err
	}

	ch := &model.Channel{
		ID:      cc.node.Generate().Int64(),
		Type:    model.ChannelTypeGuild,
		GuildID: &guildID,
		Name:    &body.ChannelName,
	}

	ctx := c.Request().Context()

	_, err = cc.db.NewInsert().Model(ch).Exec(ctx)
	if err != nil {
		log.Println("bun: insert failed:", err)
		return c.Blob(http.StatusInternalServerError, "application/json", resServerError)
	}

	dtoch := dto.GuildChannelResponse{
		ID:   ch.ID,
		Name: *ch.Name,
	}

	return c.JSON(http.StatusOK, dto.Response{
		Ok:   true,
		Data: dtoch,
	})
}

func (cc *ChannelController) CreateGuild(c *echo.Context) error {
	user := c.Get("user").(*model.User)
	body := new(dto.CreateGuildBody)

	guild := &model.Guild{
		ID:      cc.node.Generate().Int64(),
		Name:    body.GuildName,
		OwnerID: user.ID,
	}

	ctx := c.Request().Context()

	_, err := cc.db.NewInsert().Model(guild).Returning("*").Exec(ctx)
	if err != nil {
		log.Println("bun: insert failed:", err)
		return c.Blob(http.StatusInternalServerError, "application/json", resServerError)
	}

	gdto := dto.GuildResponse{
		ID:        guild.ID,
		Name:      guild.Name,
		CreatedAt: guild.CreatedAt,
	}

	return c.JSON(http.StatusOK, dto.Response{
		Ok:   true,
		Data: gdto,
	})
}
