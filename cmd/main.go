package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lamrin13/file-transfer/utils"
	"github.com/pion/webrtc/v2"
)

func main() {

	s := webrtc.SettingEngine{}
	s.DetachDataChannels()

	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	pc, err := api.NewPeerConnection(config)
	if err != nil {
		fmt.Printf("error establishing new peer connection: %s", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Send file (S) or Receive file (R): ")
	offerType, _ := reader.ReadString('\n')

	if offerType == "S\r\n" {
		utils.CreateOffer(pc, reader)
	} else {
		utils.CreateAnswer(pc, reader)
	}

	select {}
}
