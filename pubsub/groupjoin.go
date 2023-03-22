package pubsub

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"groupschemepoc1/group"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

var buff []byte
var PauseCLI bool = false

// this function's input from SDI clashes with clihanlder goroutine's input from SDI
// causing it to behave strangely
// needs to be solved
func (jr *JoinRequest) handleJoinRequest(from peer.ID) {
	fmt.Println("Join request being processed")
	if group.CurrentGroupShareKey.GroupName != jr.GroupName {
		fmt.Println("You are not part of the specified group")
		fmt.Println("[ABORT] - Join request aborted")
		return
	}
	fmt.Println("------------------------------------------")
	fmt.Println("FROM : " + from.Pretty())
	fmt.Println("GROUP : " + jr.GroupName)
	fmt.Println("MESSAGE > " + jr.Message)
	fmt.Println("------------------------------------------")
	fmt.Println("[<y> to accept/ anyhthing else to reject]")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR] - during reading the choice")
		return
	}
	PauseCLI = false
	escapeSeqLen := 0
	if runtime.GOOS == "windows" {
		escapeSeqLen = 2
	} else {
		escapeSeqLen = 1
	}
	choice = choice[0 : len(choice)-escapeSeqLen]
	fmt.Println("You chose " + choice)
	switch choice {
	case "<y>":
		fmt.Println("You chose to accept the incoming group joining request")
		joinrequestreply := &JoinRequestReply{
			GroupName: jr.GroupName,
			Host:      jr.Host,
			To:        from,
			Message:   jr.Message,
			Granted:   true,
			Key:       *group.CurrentGroupShareKey,
		}
		jrrbytes, err := json.Marshal(joinrequestreply)
		if err != nil {
			fmt.Println("[ERROR - during marshalling the group share key]")
		}

		str, err := group.MentorInfoObj.Host.NewStream(*group.MentorInfoObj.MentCTX, from, GroupJoinReplyProtocol)
		if err != nil {
			fmt.Println("[ERROR] - during establishing a stream to " + from.Pretty())
		}
		fmt.Println("[SUCCESS] - in establishing a stream to " + from.Pretty())
		str.Write(jrrbytes)
		fmt.Println("[DONE] - in sending the group key bytes")
		str.Close()
		break
	default:
		fmt.Println("You chose to reject the incoming group joinig request")
		break
	}

}

func HandleStreamJoinRequest(str network.Stream) {
	PauseCLI = true // this stream handler's input coincides with the clihanlder goroutine
	fmt.Println("Request recieved")
	go grpJoinRequest(str)
}

func grpJoinRequest(str network.Stream) {
	buff = nil
	fmt.Println("Group join request identified")
	streamReader := bufio.NewReader(str)
	buffer := make([]byte, 1)

	for {
		_, err := streamReader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Fully Recieved")
				break
			}
			fmt.Println("[ERROR] - during reading from the input stream")
			fmt.Println(err.Error())
			fmt.Println("[ABORT] - handling join request from " + str.Conn().RemotePeer())
			return
		}
		buff = append(buff, buffer...)

	}
	fmt.Println("Exited the loop")
	joinrequest := &JoinRequest{}
	err := json.Unmarshal(buff, joinrequest)
	if err != nil {
		fmt.Println("Error during unmarshalling")
		return
	}
	fmt.Println(joinrequest)
	joinrequest.handleJoinRequest(str.Conn().RemotePeer())
}

func GroupJoinRequest(ctx context.Context, host host.Host, groupName string, grphostid peer.ID, message string) {
	joinRequest := JoinRequest{
		GroupName: groupName,
		Host:      grphostid,
		Message:   message,
	}
	jrbytes, err := json.Marshal(joinRequest)
	if err != nil {
		fmt.Println("Error during marshalling the join request")
		return
	}
	fmt.Println(len(jrbytes))
	str, err := host.NewStream(ctx, grphostid, GroupJoinRequestProtocol)
	if err != nil {
		fmt.Println("Error while tyring to establish a stream to " + grphostid.Pretty())
		return
	}
	fmt.Println("[Sucessful] - in establishing a stream to [" + str.Conn().RemotePeer().Pretty() + "] ")
	fmt.Println(str.Protocol())
	fmt.Println("[Processign] - about to write following bytes to stream")
	fmt.Println(jrbytes)
	_, err = str.Write(jrbytes)
	if err != nil {
		fmt.Println("Error during sending the group join request")
		fmt.Println("Group join procedure aborted")
		return
	}
	str.Close()

}

func Test(str network.Stream) {
	fmt.Println("Test successful")

	rw := bufio.NewReadWriter(bufio.NewReader(str), bufio.NewWriter(str))

	go fmt.Print(rw.ReadByte())
}

func RecieveReply(str network.Stream) {
	buff = nil
	fmt.Println("Response recieved")
	streamReader := bufio.NewReader(str)
	buffer := make([]byte, 1)
	for {
		_, err := streamReader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Fully Recieved")
				break
			}
			fmt.Println("[ERROR] - during reading from the input stream")
			fmt.Println(err.Error())
			fmt.Println("[ABORT] - handling join response")
			return
		}
		buff = append(buff, buffer...)

	}
	fmt.Println("Exited the loop")
	joinrequestreply := &JoinRequestReply{}
	err := json.Unmarshal(buff, joinrequestreply)
	if err != nil {
		fmt.Println("[ERROR] - during unmarhsalling join request reply recieved")
		fmt.Println(buff)
	}
	fmt.Println(joinrequestreply.Key.GroupName)
	fmt.Println(joinrequestreply.Granted)
	oldgrouproom := CurrentGroupRoom
	oldgrouproom.ExitRoom()
	newgrouproom, err := JoinGroup(CurrentGroupRoom.HostP2P, CurrentGroupRoom.UserName, joinrequestreply.Key.GenerateGroupKey())
	CurrentGroupRoom = newgrouproom
	if err != nil {
		fmt.Println("Error while joining the new group")
	}
	fmt.Print("Waiting for queues to adapt")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Print(".")

}
