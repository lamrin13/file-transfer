package utils

import (
	"bufio"
	"fmt"
	"os"

	b64 "encoding/base64"

	"github.com/pion/webrtc/v2"
)

func CreateOffer(pc *webrtc.PeerConnection, reader *bufio.Reader) {

	var (
		sig  chan int = make(chan int)
		sync chan int = make(chan int)
	)

	ordered := true
	maxRetransmits := uint16(5)
	options := &webrtc.DataChannelInit{
		Ordered:        &ordered,
		MaxRetransmits: &maxRetransmits,
	}
	sender, err := pc.CreateDataChannel("Sender", options)
	if err != nil {
		fmt.Printf("error creating new data channel: %s", err)
		os.Exit(1)
	}

	pc.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		fmt.Printf("Connection status: %s", pcs.String())
		if pcs.String() == "failed" {
			fmt.Println("Connection lost...")
		}
	})

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		fmt.Printf("error creating offer: %s", err)
		os.Exit(1)
	}

	err = pc.SetLocalDescription(offer)
	if err != nil {
		fmt.Printf("error setting local description: %s", err)
		os.Exit(1)
	}
	fmt.Println(b64.StdEncoding.EncodeToString([]byte(offer.SDP)))

	fmt.Println("Enter receiver sdp: ")
	secret, _ := reader.ReadString('\n')
	sdp, _ := b64.StdEncoding.DecodeString(secret)
	answer := webrtc.SessionDescription{
		Type: 3,
		SDP:  string(sdp),
	}

	err = pc.SetRemoteDescription(answer)
	if err != nil {
		fmt.Printf("error setting remote description: %s", err)
		os.Exit(1)
	}

	sender.OnOpen(func() {
		fmt.Println("Data channel opened: ", sender.Label())

		writerReader, err := sender.Detach()
		if err != nil {
			fmt.Printf("error detach: %s", err)
		}
		go WriteSender(writerReader, sig, sync)
		go ReadSender(writerReader, sig, sync)

	})
}

func CreateAnswer(pc *webrtc.PeerConnection, reader *bufio.Reader) {
	fmt.Println("Enter sender sdp: ")
	secret, _ := reader.ReadString('\n')
	sdp, _ := b64.StdEncoding.DecodeString(secret)
	offer := webrtc.SessionDescription{
		Type: 1,
		SDP:  string(sdp),
	}

	err := pc.SetRemoteDescription(offer)
	if err != nil {
		fmt.Printf("error setting remote description: %s", err)
		os.Exit(1)
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		fmt.Printf("error creating answer: %s", err)
		os.Exit(1)
	}

	err = pc.SetLocalDescription(answer)
	if err != nil {
		fmt.Printf("error setting local description: %s", err)
		os.Exit(1)
	}

	fmt.Println(b64.StdEncoding.EncodeToString([]byte(answer.SDP)))
	pc.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		fmt.Printf("Connection status: %s\n", pcs.String())
	})

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc.OnOpen(func() {
			fmt.Println("Data channel opened ", dc.Label())
			reader, err := dc.Detach()
			fmt.Println(reader)
			if err != nil {
				fmt.Printf("Error detach: %s", err)
			}
			go ReadReceiver(reader)
		})

	})
}
