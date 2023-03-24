package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"groupschemepoc1/p2pnet"
	"math/rand"
	"time"

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

// Packet type can be the following
// <brd> - broadcast
// <brdreply> - broadcastreply
// <chat> - chat message
type packet struct {
	PacketType string
	Content    []byte
}

type BroadCastMessage struct {
	PeerId peer.ID
}

type BroadCastReplyMessage struct {
	PeerId peer.ID
}

type GroupRoom struct {
	HostP2P   *p2pnet.P2P
	GroupName string
	UserName  string
	Inbound   chan packet
	Outbound  chan packet
	State     int // 1= active and 0- dead
	SelfId    peer.ID
	psctx     context.Context
	pscancel  context.CancelFunc
	pstopic   *pbsb.Topic
	psub      *pbsb.Subscription
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
		Inbound:   make(chan packet),
		Outbound:  make(chan packet),
		SelfId:    hostp2p.Host.ID(),
		psctx:     pubSubCtx,
		State:     1,
		pscancel:  cancel,
		pstopic:   topic,
		psub:      sub,
	}
	// fmt.Println("Starting new subloop")
	go groupRoom.SubLoop()
	// fmt.Println("Starting new publoop")
	go groupRoom.PubLoop()
	fmt.Println("Started the broadcast handler for " + groupRoom.GroupName)
	go groupRoom.BroadCastHandler()
	return groupRoom, nil

}

func (gr *GroupRoom) PubLoop() {
	for {
		// fmt.Println("publoop running")
		select {
		case <-gr.psctx.Done():
			fmt.Println("PubLoop Exit")
			return
		case outpacket := <-gr.Outbound:
			// fmt.Println("Outbound message is being processed")

			packetbytes, err := json.Marshal(outpacket)
			if err != nil {
				fmt.Println("[ERROR] - during marhsalling outgoing packet")
				continue
			}

			err = gr.pstopic.Publish(gr.psctx, packetbytes)
			if err != nil {
				fmt.Println("[ERROR] - during publishing the packet to group")
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

			inpacket := &packet{}

			err = json.Unmarshal(message.Data, inpacket)
			if err != nil {
				fmt.Println("Error while unmarshalling the messages")
			}
			if message.ReceivedFrom == gr.SelfId {
				if inpacket.PacketType == "<chat>" {
					continue
				}
			}

			gr.Inbound <- *inpacket
		}
	}
}

func (gr *GroupRoom) PeerList() []peer.ID {
	return gr.pstopic.ListPeers()
}

func (gr *GroupRoom) ExitRoom() {
	defer gr.pscancel()
	endoldsession = true
	gr.State = 0
	gr.psub.Cancel()
	gr.pstopic.Close()
}

func (gr *GroupRoom) UpdateUserName(username string) {
	gr.UserName = username
}

func (gr *GroupRoom) BroadCastHandler() {
	waittime := (rand.Intn(60-20) + 20)
	fmt.Println("Broadcast wait time for this node is", waittime)
	for i := 0; i <= waittime; i++ {
		time.Sleep(1 * time.Second)
		if i == waittime {
			if broadcastrecieved {
				broadcastrecieved = false
				go gr.BroadCastHandler()
				return
			}
			if gr.State == 0 {
				fmt.Println("Ending BroadCast handler for " + gr.GroupName)
				return
			}
			broadCastMessage := &BroadCastMessage{
				PeerId: gr.SelfId,
			}
			brdbytes, err := json.Marshal(broadCastMessage)
			if err != nil {
				fmt.Println("[ERROR] - during marshalling broadcast message")
				continue
			}
			brdpacket := &packet{
				PacketType: "<brd>",
				Content:    brdbytes,
			}
			gr.Outbound <- *brdpacket
			go gr.BroadCastHandler()
		}

	}
}

func (gr *GroupRoom) BroadCastReplyHandler() {
	brdreplypacket := &BroadCastReplyMessage{
		PeerId: gr.SelfId,
	}
	brdreplybytes, err := json.Marshal(brdreplypacket)
	if err != nil {
		fmt.Println("[ERROR] - during marshalling broadcast reply")
		return
	}
	outpacket := &packet{
		PacketType: "<brdreply>",
		Content:    brdreplybytes,
	}
	gr.Outbound <- *outpacket

}
