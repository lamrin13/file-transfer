# file-transfer

Simple p2p file transfer utility using [pion webrtc](https://github.com/pion)

## How to use:

1. Run the main file using (TODO; convert to CLI tool)
```bash
go run ./cmd/main.go
```

2. Input "S" or "R" based on Sender and Receiver

3. On sender, the SDP will be printed on console; share the SDP with the Receiver

4. On receiver, provide the sender SDP; the receiver SDP will be printed on console; share the SDP with Sender

5. Provide the path of file to send on Sender machine (TODO: currently works for windows path only, generalize for other OS)

6. On receiver, accept the file being sent by typing "Y" or "N"


## TODO

- Convert to CLI
- Improve sharing SDP b/w sender and receiver
- Removing hardcoding for the windows file system
