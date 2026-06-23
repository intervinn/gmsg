package ws

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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

	e.GET("/", c.OnWS)

	return c
}

func (ctl *Controller) OnWS(c *echo.Context) error {
	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := &Client{
		Closed:      false,
		Conn:        conn,
		State:       StateHandshake,
		UserID:      0,
		ActiveGuild: 0,
	}
	ctl.cr.Add(client)

	for {
		data, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			ctl.cr.Delete(client)
			return err
		}

		if client.Closed || op == ws.OpClose {
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
	if c.State == StateHandshake {
		switch msg.Type {
		case "authenticate":
			{
				var token string
				if err := json.Unmarshal(msg.Data, &token); err != nil {
					log.Println("invalid format:", err)
					c.Update(func(c *Client) {
						c.Closed = true
					})
					return
				}

				t := strings.TrimPrefix(token, "Bearer ")
				tok, err := ctl.ts.Verify(t)
				if err != nil {
					c.Update(func(c *Client) {
						c.Closed = true
					})
					return
				}

				idstr, err := tok.Claims.GetSubject()
				if err != nil {
					log.Println("missing subject:", err)
					c.Update(func(c *Client) {
						c.Closed = true
					})
					return
				}

				id, err := strconv.ParseInt(idstr, 10, 64)
				if err != nil {
					c.Update(func(c *Client) {
						c.Closed = true
					})
					return
				}

				c.Update(func(c *Client) {
					c.UserID = id
				})
			}
		}

		return
	}

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
