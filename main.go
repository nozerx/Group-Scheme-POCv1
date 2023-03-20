package main

import (
	"fmt"
	"groupschemepoc1/p2pnet"

	"groupschemepoc1/pubsub"
)

const topic = "rex/groupscheme/test"
const service = "rex/service/groupschema/test"
const subtopic = "rex/groupscheme/test/sub"

func main() {
	ctx, host := p2pnet.EstablishP2P()
	kad_dht := p2pnet.HandleDHT(ctx, host)
	pubSub := pubsub.SetUpPubSub(ctx, host)
	p2phost := p2pnet.NewP2P(ctx, host, kad_dht, pubSub)
	grp, err := pubsub.JoinGroup(p2phost, "", topic)
	if err != nil {
		fmt.Println("Error while joining the group")
	}
	go grp.HandleInputFromSDI()
	go p2pnet.DiscoverPeers(ctx, host, kad_dht, service)
	for {

	}
}
