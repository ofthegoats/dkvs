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
	// m spawns knowing about n
	m := NewNode(
		[]string{"tcp://localhost:9999"}, "localhost",
		key,
		tcptimeout, rttperiod, fscperiod,
		9998, 5, 2,
	)
	go m.Gossip()

    time.Sleep(5 * time.Second)

	r := Rumour{
		RequestType: UpdateData,
		Key:         "newkey",
		NewValue:    "newvalue",
	}
	go m.Send("tcp://localhost:9998", r)

	for {
        log.Printf("%s:%d: %v, %v\n", n.LANIP, n.Port, n.Data, n.Neighbours)
        log.Printf("%s:%d: %v, %v\n", m.LANIP, m.Port, m.Data, m.Neighbours)
		time.Sleep(5 * time.Second)
	}
}
