package ws

import (
	"net"
	"sync"
)

type ClientState int

const (
	StateHandshake ClientState = iota
	StateReady
)

type Client struct {
	mu          sync.Mutex
	Closed      bool
	Conn        net.Conn
	State       ClientState
	UserID      int64
	ActiveGuild int64
}

func (c *Client) Update(f func(c *Client)) {
	c.mu.Lock()
	f(c)
	c.mu.Unlock()
}

type ClientRegistry struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
}

func (cr *ClientRegistry) Add(c *Client) {
	cr.mu.Lock()
	cr.clients[c] = struct{}{}
	cr.mu.Unlock()
}

func (cr *ClientRegistry) Delete(c *Client) {
	cr.mu.Lock()
	delete(cr.clients, c)
	cr.mu.Unlock()
}

func (cr *ClientRegistry) Each(f func(c *Client)) {
	cr.mu.RLock()
	for c := range cr.clients {
		f(c)
	}
	cr.mu.RUnlock()
}
