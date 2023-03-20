package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"groupschemepoc1/p2pnet"

	pbsb "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

const defaultroom = "lobby"
const defaultusername = "rex"

type chatmessage struct {
	Message    string
	SenderID   peer.ID
	SenderName string
}

type GroupRoom struct {
	HostP2P   *p2pnet.P2P
	GroupName string
	UserName  string
	Inbound   chan chatmessage
	Outbound  chan string

	selfId   peer.ID
	psctx    context.Context
	pscancel context.CancelFunc
	pstopic  *pbsb.Topic
	psub     *pbsb.Subscription
}

func SetUpPubSub(ctx context.Context, host host.Host) *pbsb.PubSub {
	pubSubHandler, err := pbsb.NewGossipSub(ctx, host)
	if err != nil {
		fmt.Println("Error during creating a new pubsub handler")
	}
	return pubSubHandler
}

func JoinGroup(hostp2p *p2pnet.P2P, username string, groupname string) (*GroupRoom, error) {
	defer fmt.Println("Joined the new group " + groupname)
	if groupname == "" {
		groupname = defaultroom
	}
	topic, err := hostp2p.PubSub.Join(groupname)
	if err != nil {
		return nil, err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		fmt.Println("Error during subscribing to new group")
		return nil, err
	}
	if username == "" {
		username = defaultusername
	}
	pubSubCtx, cancel := context.WithCancel(context.Background())

	groupRoom := &GroupRoom{
		HostP2P:   hostp2p,
		GroupName: groupname,
		UserName:  username,
		Inbound:   make(chan chatmessage),
		Outbound:  make(chan string),
		selfId:    hostp2p.Host.ID(),
		psctx:     pubSubCtx,
		pscancel:  cancel,
		pstopic:   topic,
		psub:      sub,
	}
	fmt.Println("Starting new subloop")
	go groupRoom.SubLoop()
	fmt.Println("Starting new publoop")
	go groupRoom.PubLoop()

	return groupRoom, nil

}

func (gr *GroupRoom) PubLoop() {
	for {
		fmt.Println("publoop running")
		select {
		case <-gr.psctx.Done():
			fmt.Println("PubLoop Exit")
			return
		case message := <-gr.Outbound:
			fmt.Println("Outbound message is being processed")
			m := chatmessage{
				Message:    message,
				SenderID:   gr.selfId,
				SenderName: gr.UserName,
			}
			messagebytes, err := json.Marshal(m)
			if err != nil {
				fmt.Println("Error during marshalling the message")
				continue
			}

			err = gr.pstopic.Publish(gr.psctx, messagebytes)
			if err != nil {
				fmt.Println("Error during publishing the message to group")
				continue
			}
		}
	}
}

func (gr *GroupRoom) SubLoop() {
	go gr.DisplayMessage()

	for {
		select {
		case <-gr.psctx.Done():
			break
		default:
			message, err := gr.psub.Next(gr.psctx)
			if err != nil {
				close(gr.Inbound)
				return
			}

			if message.ReceivedFrom == gr.selfId {
				fmt.Println("Message from self identified")
				continue
			}

			cm := &chatmessage{}

			err = json.Unmarshal(message.Data, cm)
			if err != nil {
				fmt.Println("Error while unmarshalling the messages")
			}

			gr.Inbound <- *cm
		}
	}
}

func (gr *GroupRoom) PeerList() []peer.ID {
	return gr.pstopic.ListPeers()
}

func (gr *GroupRoom) ExitRoom() {
	defer gr.pscancel()

	gr.psub.Cancel()
	gr.pstopic.Close()
}

func (gr *GroupRoom) UpdateUserName(username string) {
	gr.UserName = username
}
