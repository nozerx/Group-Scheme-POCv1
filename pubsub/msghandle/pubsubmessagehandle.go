package msghandle

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Chatmessage struct {
	Messagecontent string
	Messagefrom    peer.ID
	Authorname     string
}

type Packet struct {
	Type         string
	InnerContent []byte
}

func composeMessage(msg string, host host.Host) *Chatmessage {
	return &Chatmessage{
		Messagecontent: msg,
		Messagefrom:    host.ID(),
		Authorname:     host.ID().ShortString(),
	}
}

func handleInputFromSubscription(ctx context.Context, host host.Host, sub *p2ppubsub.Subscription) {
	inputPacket := &Packet{}
	for {
		inputMsg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Error while recieving next message from the subscription")
		} else {
			err := json.Unmarshal(inputMsg.Data, inputPacket)
			if err != nil {
				fmt.Println("Error while unmarshalling at packet level")
			} else {
				if inputPacket.Type == "chat" {
					chatMsg := &Chatmessage{}
					err := json.Unmarshal(inputPacket.InnerContent, chatMsg)
					if err != nil {
						fmt.Println("Error while unmarshalling at chat message level")
					} else {
						fmt.Println("[ BY-> ", inputMsg.ReceivedFrom.Pretty()[len(inputMsg.ReceivedFrom.Pretty())-6:len(inputMsg.ReceivedFrom.Pretty())], "->", chatMsg.Messagecontent)
					}
				}
			}
		}
	}
}

func handleInputFromSDI(ctx context.Context, host host.Host, topic *p2ppubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error during reading from standard input")
		} else {

			fmt.Println("Tag-> <c>")
			writeToSubscription(ctx, host, input, topic)
		}
	}
}

func writeToSubscription(ctx context.Context, host host.Host, message string, topic *p2ppubsub.Topic) {
	chatMessage := composeMessage(message, host)
	inputContent, err := json.Marshal(chatMessage)
	if err != nil {
		fmt.Println("Error while marshalling at chatmessage level")
	} else {
		packetContent := &Packet{
			InnerContent: inputContent,
			Type:         "chat",
		}
		packet, err := json.Marshal(packetContent)
		if err != nil {
			fmt.Println("Error while marshalling at packet level")
		} else {
			topic.Publish(ctx, packet)
		}
	}
}

func HandlePubSubMessages(ctx context.Context, host host.Host, sub *p2ppubsub.Subscription, topic *p2ppubsub.Topic) {
	go handleInputFromSubscription(ctx, host, sub)
	handleInputFromSDI(ctx, host, topic)
}
