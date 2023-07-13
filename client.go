package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/quic-go/quic-go"
)

const addr = "CHANGEME:4242"

func main() {
	log.Fatal(clientMain())
}

func clientMain() error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err != nil {
		return err
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	// Create a payload that normally will fit on a network interface with a 1500 MTU
	message := make([]byte, 1352)
	for i := 0; i < 1352; i++ {
		message[i] = 0xef
	}

	var input string
	for i := 0; i < 3; i++ {
		// We wait before sending to allow for the MTU to be adjusted inbetween Writes
		fmt.Scanln(&input)
		stream.SetWriteDeadline(time.Now().Add(5 * time.Second))
		fmt.Printf("Client: Sending message of size '%d' bytes\n", len(message))
		_, err = stream.Write(message)
		if err != nil {
			return err
		}

		read := 0
		for read < len(message) {
			stream.SetReadDeadline(time.Now().Add(15 * time.Second))
			buf := make([]byte, len(message))
			n, err := stream.Read(buf)
			if err != nil {
				return err
			}
			fmt.Printf("Client: Got '%d' bytes\n", n)
			read += n
		}
	}

	return nil
}
