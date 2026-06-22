package dto

import (
	"time"

	"github.com/intervinn/gmsg/api/model"
)

type CreateMessageBody struct {
	Content string `json:"content"`
}

type CreateChannelBody struct {
	ChannelName string `json:"channel_name"`
}

type EventMessage struct {
	Type string         `json:"t"`
	Data *model.Message `json:"d"`
}

type GuildChannelResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateGuildBody struct {
	GuildName string `json:"guild_name"`
}

type GuildResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
