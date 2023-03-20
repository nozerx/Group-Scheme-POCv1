package main

import (
	"groupschemepoc1/p2pnet"
	"groupschemepoc1/pubsub"
	"groupschemepoc1/pubsub/msghandle"
)

const topic = "rex/groupscheme/test"
const service = "rex/service/groupschema/test"

func main() {
	ctx, host := p2pnet.EstablishP2P()
	kad_dht := p2pnet.HandleDHT(ctx, host)
	sub, top := pubsub.HandlePubSub(ctx, host, topic)
	go p2pnet.DiscoverPeers(ctx, host, kad_dht, service)
	msghandle.HandlePubSubMessages(ctx, host, sub, top)
}
