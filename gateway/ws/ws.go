package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/intervinn/gmsg/gateway/dto"
	"github.com/intervinn/gmsg/gateway/service"
	"github.com/labstack/echo/v5"
	"github.com/nats-io/nats.go"
)

type Message struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

type Controller struct {
	e  *echo.Echo
	nc *nats.Conn
	ts *service.TokenService

	cr *ClientRegistry
	gr *GuildRegistry
}

func New(e *echo.Echo, nc *nats.Conn, ts *service.TokenService) *Controller {
	c := &Controller{
		e:  e,
		nc: nc,
		ts: ts,
	}

	c.cr = &ClientRegistry{
		clients: map[*Client]struct{}{},
	}

	c.gr = &GuildRegistry{
		cr:   c.cr,
		nc:   c.nc,
		subs: map[int64]*nats.Subscription{},
	}

	e.GET("/ws", c.OnWS)

	return c
}

func (ctl *Controller) OnWS(c *echo.Context) error {
	tokenstr := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(tokenstr, "Bearer ")

	_, err := ctl.ts.Verify(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.Response{
			Ok:    false,
			Error: "Unauthorized",
		})
	}

	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
	if err != nil {
		return err
	}

	client := &Client{
		Conn:        conn,
		State:       StateReady,
		ActiveGuild: 0,
	}
	ctl.cr.Add(client)

	for {
		data, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			ctl.cr.Delete(client)
			return err
		}

		if op == ws.OpClose {
			ctl.cr.Delete(client)
			break
		}

		msg := new(Message)
		if err := json.Unmarshal(data, msg); err != nil {
			log.Println("failed to unmarshal:", err)
			continue
		}

		ctl.OnMsg(client, msg)
	}

	return nil
}

func (ctl *Controller) OnMsg(c *Client, msg *Message) {
	switch msg.Type {
	case "set_active":
		{
			var id int64
			if err := json.Unmarshal(msg.Data, &id); err != nil {
				log.Println("invalid format:", err)
				return
			}

			log.Println("updating active guild:", id)
			ctl.gr.Refresh()
			c.Update(func(c *Client) {
				c.ActiveGuild = id
			})
		}
	}
}
