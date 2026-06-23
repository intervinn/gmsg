package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gobwas/ws/wsutil"
	"github.com/nats-io/nats.go"
)

type GuildMessage struct {
	Type string `json:"t"`
	Data struct {
		ID        int64     `json:"id"`
		Content   string    `json:"content"`
		AuthorID  int64     `json:"author_id"`
		GuildID   int64     `json:"guild_id"`
		ChannelID int64     `json:"channel_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"d"`
}

type GuildRegistry struct {
	mu sync.RWMutex
	cr *ClientRegistry
	nc *nats.Conn

	subs map[int64]*nats.Subscription
}

func (gr *GuildRegistry) Refresh() error {
	gr.mu.Lock()
	defer gr.mu.Unlock()

	ids := map[int64]struct{}{}
	gr.cr.Each(func(c *Client) {
		ids[c.ActiveGuild] = struct{}{}
	})

	for id := range gr.subs {
		if _, ok := ids[id]; !ok {
			gr.subs[id].Unsubscribe()
			delete(gr.subs, id)
		}
	}

	for id := range ids {
		if _, ok := gr.subs[id]; ok {
			continue
		}

		sub, err := gr.nc.Subscribe(fmt.Sprintf("guild.%v.message", id), gr.onMsg)
		if err != nil {
			log.Println("failed to subscribe:", err)
			continue
		}

		gr.subs[id] = sub
	}

	return nil
}

func (gr *GuildRegistry) onMsg(m *nats.Msg) {
	log.Println("received nats message")
	msg := new(GuildMessage)
	data := m.Data
	if err := json.Unmarshal(data, msg); err != nil {
		log.Println("failed to unmarshal:", err)
	}
	id := msg.Data.GuildID

	gr.cr.Each(func(c *Client) {
		if c.ActiveGuild == id && c.UserID != msg.Data.AuthorID {
			wsutil.WriteServerText(c.Conn, data)
		}
	})
}
