# wsutil

> Implements websocket communication

## Usage example

```go
// Package ws implements websocket communication
package ws

import (
	"github.com/avakarev/go-util/wsutil"
	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
)

// Peer implements connection's peer
type Peer struct {
	wsutil.BasePeer
	id string
}

// Match checks whether given peer is quialified to recevive email
func (p *Peer) Match(id string) bool {
	return p.id == id
}

// Hub maintains the set of active peers and handles communication
type Hub struct {
	wsutil.BroadcastHub
}

// Handler handles websocket requests from the peer
func (h *Hub) Handler(conn *websocket.Conn) {
	// when the function returns, unregister the peer and close the connection
	defer func() {
		h.Unregister(conn)
	}()

	// register the peer
	h.Register(&wsutil.PeerRequest{Conn: conn, Peer: &Peer{id: conn.Params("id")}})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warn().Err(err).Msg("ws read error")
			}
			// causes call of the deferred function, i.e. closes the connection on error
			return
		}
	}
}

// NewHub returns new hub value
func NewHub() *Hub {
	return &Hub{BroadcastHub: *wsutil.NewBroadcastHub()}
}
```
