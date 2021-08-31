package main

import (
	"crypto/rand"
	"io"
	"log"
	"time"
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
        log.Printf("%s:%d: %v, %v\n", n.LANIP, n.Port, n.Data, n.Neighbours)
		time.Sleep(5 * time.Second)
	}
}
