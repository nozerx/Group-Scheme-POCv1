package pubsub

import "github.com/libp2p/go-libp2p/core/peer"

type ServicePeer struct {
	Id     int
	PeerId peer.ID
}
