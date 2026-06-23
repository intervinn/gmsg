package dto

import "time"

type RegisterBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResult struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Authorization string `json:"authorization"`
}

type LoginResult struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Authorization string `json:"authorization"`
}

type PublicUser struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
}
