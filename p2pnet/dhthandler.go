package p2pnet

import (
	"context"
	"fmt"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

func initDHT(ctx context.Context, host host.Host) *dht.IpfsDHT {
	kad_dht, err := dht.New(ctx, host)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error while creating the DHT")
	} else {
		fmt.Println("Successful in creating a new DHT")
	}
	return kad_dht
}

func bootstrapDHT(ctx context.Context, host host.Host, kad_dht *dht.IpfsDHT) {
	err := kad_dht.Bootstrap(ctx)
	if err != nil {
		fmt.Println("Error while setting the dht into bootstrap state")
	} else {
		fmt.Println("Successful in setting the DHT into bootstrap state")
	}

	for _, bootstrapPeer := range dht.GetDefaultBootstrapPeerAddrInfos() {
		err := host.Connect(ctx, bootstrapPeer)
		if err != nil {
			fmt.Println("Error while connecting to default bootstrap peer : ", bootstrapPeer.ID)
		} else {
			fmt.Println("Successfully connected to bootpeer : ", bootstrapPeer.ID)
		}
	}
	fmt.Println("Done with all connections to default bootstrap peers")
}

func HandleDHT(ctx context.Context, host host.Host) *dht.IpfsDHT {
	kad_dht := initDHT(ctx, host)
	bootstrapDHT(ctx, host, kad_dht)
	return kad_dht
}
