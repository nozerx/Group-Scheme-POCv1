package pubsub

import (
	"context"
	"fmt"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

func newPubSubService(ctx context.Context, host host.Host) *p2ppubsub.PubSub {
	pubSubService, err := p2ppubsub.NewGossipSub(ctx, host)
	if err != nil {
		fmt.Println("Error while establishing a pubsub service")
	} else {
		fmt.Println("Successful in establishing a pubsub service")
	}
	return pubSubService
}

func HandlePubSub(ctx context.Context, host host.Host, topic string) (*p2ppubsub.Subscription, *p2ppubsub.Topic) {
	pubSubServie := newPubSubService(ctx, host)
	pubSubTopic, err := pubSubServie.Join(topic)
	if err != nil {
		fmt.Println("Error while joining the topic", topic)
	}

	pubSubSubscription, err := pubSubServie.Subscribe(topic)
	if err != nil {
		fmt.Println("Error while subscribing to topic", topic)
	}
	return pubSubSubscription, pubSubTopic
}
