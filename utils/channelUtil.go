package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadReceiver(d io.ReadWriter) {
	reader := bufio.NewReader(os.Stdin)
	var file *os.File
	for {
		var buffer []byte = make([]byte, 65535)
		n, err := d.Read(buffer)
		if err != nil {
			fmt.Printf("error while reading message: %s", err)
			os.Exit(1)
		}
		var path string
		if strings.HasPrefix(string(buffer[:n]), "Want to receive file:") {
			fmt.Printf("Message from sender: %s\n", string(buffer[:n]))
			path = strings.Split(string(buffer[:n]), ": ")[1]
			fmt.Printf("Reply with Y or N: ")
			message, _ := reader.ReadString('\n')
			if strings.TrimRight(message, "\r\n") == "Y" {
				path = strings.TrimRight(path, "\r\n")
				file, err = os.Create(path)
				if err != nil {
					fmt.Printf("cannot create %s: %s", file.Name(), err)
					os.Exit(1)
				}
				d.Write([]byte(message))
			}
		} else {
			WriteFile(d, file)
		}

	}
}

func WriteSender(d io.Writer, sig <-chan int, sync chan<- int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nPath of the file to send: ")
		path, _ := reader.ReadString('\n')
		fullpath := strings.Split(path, "\\")
		filename := fullpath[len(fullpath)-1]
		_, err := d.Write([]byte("Want to receive file: " + filename))
		if err != nil {
			panic(err)
		}
		sync <- 1
		flag := <-sig
		if flag == 1 {
			SendFile(d, path)
		} else {
			continue
		}
	}
}

func ReadSender(d io.Reader, sig chan<- int, sync <-chan int) {
	for {
		var buffer []byte = make([]byte, 65535)
		<-sync
		n, err := d.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Message from receiver: %s", string(buffer[:n]))
		if strings.TrimRight(string(buffer[:n]), "\r\n") == "Y" {
			sig <- 1
		} else {
			sig <- 0
		}
	}
}
