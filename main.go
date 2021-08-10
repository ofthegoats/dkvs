package main

import (
	"log"
	"time"
)

func main() {
	tcptimeout := 1 * time.Second
	rttperiod := 60 * time.Second
	fscperiod := 2 * time.Second // low, for testing
	n := NewNode(
		[]string{}, "tcp://localhost:9999",
		tcptimeout, rttperiod, fscperiod,
		5, 2,
	)
	go n.Gossip()
	time.Sleep(2 * time.Second) // time to catch up
	r := Rumour{
		RequestType: UpdateData,
		Key:         "newkey",
		NewValue:    "newvalue",
	}
	go n.Send("tcp://localhost:9999", r)

	// m spawns knowing about n
	m := NewNode(
		[]string{"tcp://localhost:9999"}, "tcp://localhost:9998",
		tcptimeout, rttperiod, fscperiod,
		5, 2,
	)
	go m.Gossip()
	time.Sleep(2 * time.Second) // time to catch up

	for {
		log.Printf("%s: %v, %v\n", n.socket, n.Data, n.Neighbours)
		log.Printf("%s: %v, %v\n", m.socket, m.Data, m.Neighbours)
		time.Sleep(5 * time.Second)
	}
}
