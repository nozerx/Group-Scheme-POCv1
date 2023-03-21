package pubsub

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"
)

var endoldsession bool

func (gr *GroupRoom) HandleInputFromSDI() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error during reading input from the stream")
			continue
		}
		if input[:5] == "<cmd>" {
			fmt.Println("These are the available commands")
			fmt.Println("1.Change Group\n2.Change UserName\n3.List Group Peers")
			var choice int
			fmt.Scanln(&choice)
			switch choice {
			case 1:
				fmt.Println("Enter the new group name to join it")
				groupName, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error during reading input from the SDI")
					continue
				}
				escapeSeqLen := 0
				if runtime.GOOS == "windows" {
					escapeSeqLen = 2
				} else {
					escapeSeqLen = 1
				}
				groupName = groupName[0 : len(groupName)-escapeSeqLen]
				oldGroupRoom := gr
				oldGroupRoom.ExitRoom()
				time.Sleep(2 * time.Second)
				newGroupRoom, err := JoinGroup(gr.HostP2P, gr.UserName, groupName)
				if err != nil {
					fmt.Println("Error while joining the new group")
					continue
				}
				gr = newGroupRoom
				fmt.Print("Waiting for queues to adapt")
				time.Sleep(1 * time.Second)
				fmt.Print(".")
				time.Sleep(1 * time.Second)
				fmt.Print(".")
				time.Sleep(1 * time.Second)
				fmt.Print(".")
				break
			case 2:
				fmt.Println("Enter the new user name")
				userName, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error during reading input from the SDI")
					continue
				}
				escapeSeqLen := 0
				if runtime.GOOS == "windows" {
					escapeSeqLen = 2
				} else {
					escapeSeqLen = 1
				}
				userName = userName[0 : len(userName)-escapeSeqLen]
				gr.UpdateUserName(userName)
				break
			case 3:
				fmt.Println("These are the list of peers in the Group " + gr.GroupName)
				for _, peer := range gr.PeerList() {
					fmt.Println(peer)
				}
				break
			default:
				fmt.Println("Bad command")

			}
		} else {
			escapeSeqLen := 0
			if runtime.GOOS == "windows" {
				escapeSeqLen = 2
			} else {
				escapeSeqLen = 1
			}
			msg := input[0 : len(input)-escapeSeqLen]
			// fmt.Println("sending message to outbound queue")
			go func() {
				gr.Outbound <- msg
				// fmt.Println("Message sent to outbound queue")
			}()

		}
	}
}

func (gr *GroupRoom) DisplayMessage() {
	fmt.Println("Starting DisplayMessage Loop")
	for msg := range gr.Inbound {
		select {
		case <-gr.psctx.Done():
			fmt.Println("Exiting DisplayMessage Loop")
			return
		default:
			fmt.Println("--------------------------------------------------")
			fmt.Printf("%s: %s\n", msg.SenderName, msg.Message)
			fmt.Println("--------------------------------------------------")
		}
	}
}
