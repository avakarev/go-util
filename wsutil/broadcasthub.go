package wsutil

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

// BroadcastHub maintains the set of active peers and handles communication
type BroadcastHub struct {
	isRunning bool

	// Registered peers
	peers map[Conn]Peer

	// Register requests from the peers
	register chan *PeerRequest

	// Unregister requests from peers
	unregister chan Conn

	// Broadcast message to peers that match the topic
	broadcast chan *Event
}

func (h *BroadcastHub) run() {
	for {
		select {
		case req := <-h.register:
			h.peers[req.Conn] = req.Peer
		case conn := <-h.unregister:
			// remove the peer from the hub
			delete(h.peers, conn)
			if err := conn.Close(); err != nil {
				log.Warn().Err(err).Send()
			}
		case event := <-h.broadcast:
			for conn, peer := range h.peers {
				// send to each peer in parallel so we don't block on a slow peer
				go func(conn Conn, peer Peer) {
					peer.Lock()
					defer peer.Unlock()
					if peer.IsClosing() {
						return
					}
					if !peer.Match(event.Topic) {
						return
					}
					if err := conn.WriteMessage(TextMessage, event.Data); err != nil {
						log.Error().Err(err).Msg("ws write error")
						peer.Close()
						if err := conn.WriteMessage(CloseMessage, []byte{}); err != nil {
							log.Warn().Err(err).Send()
						}
						if err := conn.Close(); err != nil {
							log.Warn().Err(err).Send()
						}
						h.unregister <- conn
					}
				}(conn, peer)
			}
		}
	}
}

// Run runs the broadcast hub
func (h *BroadcastHub) Run() {
	if h.isRunning {
		return
	}
	go h.run()
	h.isRunning = true
}

// Register handles new peer connection request
func (h *BroadcastHub) Register(req *PeerRequest) {
	h.register <- req
}

// Unregister removes given peer connection from the hub
func (h *BroadcastHub) Unregister(conn Conn) {
	h.unregister <- conn
}

// Broadcast sends given message to all connected peers
func (h *BroadcastHub) Broadcast(e *Event) {
	h.broadcast <- e
}

// BroadcastJSON marshals and sends given pointer destination to all connected peers
func (h *BroadcastHub) BroadcastJSON(topic string, v any) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	h.Broadcast(&Event{Topic: topic, Data: bytes})
	return nil
}

// TryBroadcastJSON marshals and sends given pointer destination to all connected peers
func (h *BroadcastHub) TryBroadcastJSON(topic string, v any) {
	if err := h.BroadcastJSON(topic, v); err != nil {
		log.Error().Err(err).Msgf("can't broadcast json (topic=%s)", topic)
	}
}

// NewBroadcastHub returns new hub value
func NewBroadcastHub() *BroadcastHub {
	return &BroadcastHub{
		peers:      make(map[Conn]Peer),
		register:   make(chan *PeerRequest),
		unregister: make(chan Conn),
		broadcast:  make(chan *Event),
	}
}
