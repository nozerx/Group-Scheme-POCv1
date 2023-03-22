package group

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type GroupKeyShare struct {
	GroupName string
	Host      peer.ID
	Key       string
}

type MentorInfo struct {
	PeerId  peer.ID
	Host    host.Host
	MentCTX *context.Context
}

var CurrentGroupShareKey *GroupKeyShare //doesnot hanlde the default key
var MentorInfoObj *MentorInfo
