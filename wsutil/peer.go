package wsutil

import "sync"

// Peer defines connection's peer
type Peer interface {
	// Lock locks mutex
	Lock()
	// Unlock unlocks mutex
	Unlock()
	// Close starts closing routine
	Close()
	// IsClosing checks whether peer started closing
	IsClosing() bool
	// Match checks whether given peer is qualified to receive message
	Match(topic string) bool
}

// BasePeer implements embeddable base peer
type BasePeer struct {
	mu        sync.Mutex
	isClosing bool
}

// Lock locks mutex
func (p *BasePeer) Lock() {
	p.mu.Lock()
}

// Unlock unlocks mutex
func (p *BasePeer) Unlock() {
	p.mu.Unlock()
}

// Close starts closing routine
func (p *BasePeer) Close() {
	p.isClosing = true
}

// IsClosing checks whether peer started closing
func (p *BasePeer) IsClosing() bool {
	return p.isClosing
}

// PeerRequest defines peer's connection request
type PeerRequest struct {
	Conn Conn
	Peer Peer
}
