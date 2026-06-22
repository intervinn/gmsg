package dto

type Response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Ok      bool   `json:"ok,omitempty"`
	Data    any    `json:"data,omitempty"`
}
