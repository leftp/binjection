// +build !windows

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
)

func MakePipe(pipename string) string {
	if _, err := os.Stat(pipename); os.IsNotExist(err) {
		// Create named pipe
		syscall.Mkfifo(pipename, 0600)
	} else if err != nil {
		log.Fatal(err)
	}
	return pipename
}

func ListenPipeDry(namedPipe string) {

	MakePipe(namedPipe)
	// Open named pipe for reading
	fmt.Println("Opening named pipe for reading")
	for {
		var buff bytes.Buffer
		stdout, err := os.OpenFile(namedPipe, os.O_RDONLY, 0600)
		if err != nil {
			log.Fatalf("Open(%s) failed: %v", namedPipe, err)
		}
		io.Copy(&buff, stdout)
		stdout.Close()

		go handleDryConnection(buff)
	}
}

func ListenPipeWet(namedPipe string) {

	MakePipe(namedPipe)
	// Open named pipe for writing
	fmt.Println("Opening named pipe for writing")
	for {
		if lastBytes != nil {
			stdout, err := os.OpenFile(namedPipe, os.O_WRONLY, 0600)
			if err != nil {
				log.Fatalf("Open(%s) failed: %v", namedPipe, err)
			}
			_, err = io.Copy(stdout, bytes.NewReader(lastBytes))
			stdout.Close()

			log.Println("Wrote wet bytes: ", len(lastBytes))

			if err != nil {
				log.Fatalf("Error on writing to pipe: %v", err)
			}
			lastBytes = nil
		}
	}
}

var lastBytes []byte

func handleDryConnection(buff bytes.Buffer) {

	i, err := Inject(buff.Bytes())
	if err != nil {
		log.Fatalf("Error injecting: %v", err)
	}
	log.Println("Set lastBytes: ", len(lastBytes))
	lastBytes = i
}
