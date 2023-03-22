package pubsub

import (
	"groupschemepoc1/group"

	"github.com/libp2p/go-libp2p/core/peer"
)

type ServicePeer struct {
	Id     int
	PeerId peer.ID
}

type GroupKeyShare struct {
	GroupName string
	Host      peer.ID
	Key       string
}

const GroupJoinRequestProtocol = "/rex/request"
const GroupJoinReplyProtocol = "rex/reply"

type JoinRequest struct {
	GroupName string
	Host      peer.ID
	Message   string
}

type JoinRequestReply struct {
	GroupName string
	Host      peer.ID
	To        peer.ID
	Message   string
	Granted   bool
	Key       group.GroupKeyShare
}
