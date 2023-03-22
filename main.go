package main

import (
	"fmt"
	"groupschemepoc1/p2pnet"

	"groupschemepoc1/pubsub"

	"github.com/libp2p/go-libp2p/core/protocol"
)

const topic = "rex/groupscheme/test"
const service = "rex/service/groupschema/test"
const subtopic = "rex/groupscheme/test/sub"
const testProtocol = "/test/0.1"
const GroupJoinRequestProtocol = "/rex/request"
const GroupJoinReplyProtocol = "rex/reply"

var testsucc bool = false

func main() {
	ctx, host := p2pnet.EstablishP2P()
	host.SetStreamHandler(protocol.ID(GroupJoinRequestProtocol), pubsub.HandleStreamJoinRequest)
	host.SetStreamHandler(protocol.ID(GroupJoinReplyProtocol), pubsub.RecieveReply)
	kad_dht := p2pnet.HandleDHT(ctx, host)
	pubSub := pubsub.SetUpPubSub(ctx, host)
	p2phost := p2pnet.NewP2P(ctx, host, kad_dht, pubSub)
	grp, err := pubsub.JoinGroup(p2phost, "", topic)
	pubsub.CurrentGroupRoom = grp
	if err != nil {
		fmt.Println("Error while joining the group")
	}
	go grp.HandleInputFromSDI(ctx, host)
	go p2pnet.DiscoverPeers(ctx, host, kad_dht, service)
	for {

	}
}
