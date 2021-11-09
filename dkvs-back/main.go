package main

import (
	"crypto/rand"
	"io"
	"time"

	. "github.com/ofthegoats/dkvs/dkvs-back/node"
)

func main() {
	tcptimeout := 1 * time.Second
	rttperiod := 300 * time.Second
	fscperiod := 300 * time.Second
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)
	n := NewNode(
		[]string{}, "localhost",
		key,
		tcptimeout, rttperiod, fscperiod,
		9999, 5, 2,
	)
	go n.Gossip()
	for {
		// inelegant infinite loop, replace with workgroups TODO
	}
}
