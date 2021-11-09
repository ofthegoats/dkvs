package main

import (
	"time"

	. "github.com/ofthegoats/dkvs/dkvs-back/node"
)

func main() {
	tcptimeout := 1 * time.Second
	rttperiod := 300 * time.Second
	fscperiod := 300 * time.Second
	key := []byte("thishasthirtytwobytesasdfasdfasd")
	n := NewNode(
		[]string{}, "localhost",
		key,
		tcptimeout, rttperiod, fscperiod,
		9999, 5, 2,
	)
	go n.Gossip()
	n.UpdateData("key", "val") // test data
	for {
		// inelegant infinite loop, replace with workgroups TODO
	}
}
