package main

import (
	"crypto/rand"
	"io"
	"log"
	"time"
)

func main() {
	tcptimeout := 1 * time.Second
	rttperiod := 60 * time.Second
	fscperiod := 2 * time.Second // low, for testing
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
        log.Printf("%s:%d: %v, %v\n", n.LANIP, n.Port, n.Data, n.Neighbours)
		time.Sleep(5 * time.Second)
	}
}
