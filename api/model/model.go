package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ChannelType int16

const (
	ChannelTypeGuild  = 2
	ChannelTypeDirect = 4
	ChannelTypeGroup  = 6
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64     `bun:"id,pk"`
	Username      string    `bun:"username,unique,notnull,type:text"`
	PasswordHash  string    `bun:"password_hash,notnull,type:text"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
}

type Channel struct {
	bun.BaseModel `bun:"table:channels,alias:c"`
	ID            int64       `bun:"id,pk"`
	Name          *string     `bun:"name,type:text"`
	Type          ChannelType `bun:"type,notnull"`
	GuildID       *int64      `bun:"guild_id"`
	Guild         *Guild      `bun:"rel:belongs-to,join:guild_id=id"`
}

type Guild struct {
	bun.BaseModel `bun:"table:guilds,alias:g"`
	ID            int64     `bun:"id,pk"`
	Name          string    `bun:"name,notnull"`
	OwnerID       int64     `bun:"owner_id,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`

	Owner    *User      `bun:"rel:belongs-to,join:owner_id=id"`
	Channels []*Channel `bun:"rel:has-many,join:id=guild_id"`
}

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m" json:"-"`
	ID            int64     `bun:"id,pk" json:"id"`
	Content       string    `bun:"content,notnull" json:"content"`
	AuthorID      int64     `bun:"author_id,notnull" json:"author_id"`
	GuildID       *int64    `bun:"guild_id" json:"guild_id"`
	ChannelID     int64     `bun:"channel_id,notnull" json:"channel_id"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp" json:"created_at"`

	Author  *User    `bun:"rel:belongs-to,join:author_id=id" json:"author,omitempty"`
	Guild   *Guild   `bun:"rel:belongs-to,join:guild_id=id" json:"guild,omitempty"`
	Channel *Channel `bun:"rel:belongs-to,join:channel_id=id" json:"channel,omitempty"`
}

type GuildMember struct {
	bun.BaseModel `bun:"table:guild_members,alias:gm"`

	GuildID int64 `bun:"guild_id,pk"`
	UserID  int64 `bun:"user_id,pk"`

	JoinedAt time.Time `bun:"joined_at,nullzero"`

	Guild *Guild `bun:"rel:belongs-to,join:guild_id=id"`
	User  *User  `bun:"rel:belongs-to,join:user_id=id"`
}
