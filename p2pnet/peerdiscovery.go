package p2pnet

import (
	"context"
	"fmt"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func DiscoverPeers(ctx context.Context, host host.Host, kad_dht *dht.IpfsDHT, service string) {
	routingDiscovery := routing.NewRoutingDiscovery(kad_dht)
	util.Advertise(ctx, routingDiscovery, service)
	fmt.Println("Successfull in advertising the service")
	for {
		peerChannel, err := routingDiscovery.FindPeers(ctx, service)
		if err != nil {
			fmt.Println("Error while finding peers with same service")
		} else {
			fmt.Println("Successfull in finding some peers")
		}

		for peerAddr := range peerChannel {
			err := host.Connect(ctx, peerAddr)
			if err != nil {
				// fmt.Println("Error while connecting to peer :", peerAddr.ID)
			} else {
				fmt.Println("Successful in connecting to peer :", peerAddr.ID)
			}
		}
	}

}
