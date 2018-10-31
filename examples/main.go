package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"cherrychain/p2p"

	logging "cherrychain/clogging"

	protocol "github.com/libp2p/go-libp2p-protocol"
)

var (
	bootstrapPeers = []string{}
	log            = logging.MustGetLogger("MAIN")
)

const (
	ip         = "0.0.0.0"
	protocolID = "/cherryCahin/1.0"
	networkID  = "cherry-test"
)

func main() {
	port := flag.Int("sp", 3000, "listen port")
	dest := flag.String("d", "", "Dest MultiAddr String")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	p2pNetwork := p2p.New(ctx, fmt.Sprintf("/ip4/%s/tcp/%d", ip, *port))

	if err := p2pNetwork.StartSysEventLoop(ctx); err != nil {
		cancel()
	}

	if *dest != "" {
		bootstrapPeers = append(bootstrapPeers, *dest)
	}
	//  else {
	// 	maddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, *port))
	// 	info, _ := peerstore.InfoFromP2pAddr(maddr)
	// 	p2pNetwork.Host.NewStream(context.Background(), info.ID, protocolID)
	// }

	pID := protocol.ID(protocolID)
	p2pNetwork.Host.SetStreamHandler(pID, p2pNetwork.HandleStream)

	log.Notice(fmt.Sprintf("./main -d /ip4/%s/tcp/%d/ipfs/%s \n", ip, *port, p2pNetwork.Host.ID().Pretty()))

	conf := p2p.Config{
		BootstrapPeers: bootstrapPeers,
		MinPeers:       -10,
		NetworkID:      networkID,
		ProtocolID:     pID,
		Notify:         p2pNetwork.Notify,
	}

	if _, err := p2pNetwork.Bootstrap(p2pNetwork, conf); err == nil {
		go writeData(p2pNetwork)
		go readData(p2pNetwork)
	}

	select {}
}

func writeData(network *p2p.P2P) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		network.Write([]byte(sendData))
	}
}

func readData(network *p2p.P2P) {
	for {
		cap := make([]byte, 1000)
		network.Read(cap)
		fmt.Printf("\x1b[32m%s\x1b[0m> ", string(cap))
	}
}