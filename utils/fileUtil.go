package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func SendFile(d io.Writer, path string) {
	file, err := os.Open(strings.TrimRight(path, "\r\n"))
	if err != nil {
		fmt.Printf("cannot open the file: %s", err)
		os.Exit(1)
	}

	defer file.Close()
	stat, _ := file.Stat()
	d.Write([]byte("Size: " + strconv.Itoa(int(stat.Size()))))
	reader := bufio.NewReader(file)

	buf := make([]byte, 65535)
	sent := 0
	since := time.Now()
	for {
		n, err := reader.Read(buf)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Printf("\rSending %d KB", sent/1024)
		_, err = d.Write(buf[:n])
		if err != nil {
			fmt.Printf("Error %s", err)
		}
		sent += n

	}
	elapsed := time.Since(since).Seconds()
	fmt.Printf("\nSent data in %f seconds at %.3f MBps", elapsed, float64(sent)/1024/1024/elapsed)
}

func WriteFile(d io.Reader, file *os.File) {
	st := time.Now()
	totalBytes := 0
	for {

		var buffer []byte = make([]byte, 65535)

		n, err := d.Read(buffer)
		if err != nil {
			fmt.Println("error reading data")
		}
		_, err = io.Copy(file, bytes.NewReader(buffer[:n]))
		if err != nil {
			fmt.Printf("error writing data %s", err)
		}
		totalBytes += n
		if n != 65535 {

			fmt.Printf("\ntotal time: %f, rate: %fMBps", time.Since(st).Seconds(), float64(totalBytes)/1024/1024/time.Since(st).Seconds())
			file.Close()
			break
		}

		fmt.Printf("\rTotal Bytes received: %d", totalBytes)

	}
	file.Close()
}
