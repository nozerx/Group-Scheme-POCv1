package streamhandler

import "github.com/libp2p/go-libp2p/core/peer"

const GroupJoinRequestProtocol = "rex/group/member/join/request"

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
	Key       string
}
