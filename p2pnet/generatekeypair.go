package p2pnet

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func GenerateKey() crypto.PrivKey {
	privkey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		fmt.Println("Error while generating the key pair")
		return nil
	} else {
		fmt.Println("Sucessful in generating the key pair")
		peerID, _ := peer.IDFromPrivateKey(privkey)
		fmt.Println("Peer ID : ", peerID)
		return privkey
	}

}
