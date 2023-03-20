package p2pnet

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

func EstablishP2P() (context.Context, host.Host) {
	prvkey := GenerateKey()
	identity := libp2p.Identity(prvkey)
	nat := libp2p.NATPortMap()
	holepunch := libp2p.EnableHolePunching()
	// mux := libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport)
	// tlstransport, err := libp2ptls.New(ID, prvkey, []upgrader.StreamMuxer{})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// security := libp2p.Security(libp2ptls.ID, tlstransport)
	transport := libp2p.Transport(tcp.NewTCPTransport)
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"), identity, nat, holepunch, transport)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Error while setting up the node")
	} else {
		fmt.Println("Successfull in setting up the node")
	}
	ctx := context.Background()
	return ctx, host

}
