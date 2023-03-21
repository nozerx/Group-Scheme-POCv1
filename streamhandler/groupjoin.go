package streamhandler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (jr *JoinRequest) handleJoinRequest() {
	fmt.Println("Join request being processed")
}

func HandleStreamJoinRequest(str network.Stream) {
	fmt.Println("Request recieved")
	go GrpJoinRequest(str)
	for {
		fmt.Println("hacked")
	}
}

func GrpJoinRequest(str network.Stream) {
	fmt.Println("Group join request identified")
	// streamReader := bufio.NewReader(str)
	// buffer := make([]byte, 100)
	// _, err := streamReader.Read(buffer)

	// if err != nil {
	// 	if err == io.EOF {
	// 		fmt.Println("Error during readinf from the input stream")
	// 		fmt.Println("ERROR : " + err.Error())
	// 		fmt.Println("Read")
	// 		fmt.Println(buffer)
	// 		fmt.Println("---------------------------------")
	// 		return
	// 	}
	// 	fmt.Println("Error during reading from the input stream")
	// 	fmt.Println("ERROR : " + err.Error())
	// 	return
	// }
	// fmt.Println("Read")
	// fmt.Println(buffer)
	// fmt.Println("---------------------------------")
	// joinrequest := &JoinRequest{}
	// err = json.Unmarshal(buffer, joinrequest)
	// if err != nil {
	// 	fmt.Println("Error during unmarshalling")
	// 	return
	// }
	// fmt.Println(joinrequest)
	// joinrequest.handleJoinRequest()
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
	// writer := bufio.NewWriter(str)
	fmt.Println("[Processign] - about to write following bytes to stream")
	fmt.Println(jrbytes)
	// // _, err = writer.Write(jrbytes)
	// if err != nil {
	// 	fmt.Println("Error during sending the group join request")
	// 	fmt.Println("Group join procedure aborted")
	// 	return
	// }

}

func Test(str network.Stream) {
	fmt.Println("Test successful")

	rw := bufio.NewReadWriter(bufio.NewReader(str), bufio.NewWriter(str))

	go fmt.Print(rw.ReadByte())
}
