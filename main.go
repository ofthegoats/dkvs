package main

import (
	"log"
	"time"
)

func main() {
	n := NewNode([]string{"tcp://127.0.0.1:9999"}, "tcp://127.0.0.1:9998", 2, 1)
	go n.Gossip()
	time.Sleep(1 * time.Second)
	m := NewNode([]string{"tcp://127.0.0.1:9998"}, "tcp://127.0.0.1:9999", 2, 1)
	go m.Gossip()
	time.Sleep(1 * time.Second)

	r1 := NewRumour("key1", "val1")
	n.Send(n.socket, r1) // send to self, hack for alpha tests
	r2 := NewRumour("key2", "val2")
	n.Send(n.socket, r2) // send to self, hack for alpha tests
	r3 := NewRumour("key1", "val3")
	n.Send(n.socket, r3) // send to self, hack for alpha tests

	time.Sleep(1 * time.Second) // give time to catch up

	log.Printf("%v\n", n.Data)
	log.Printf("%v\n", m.Data)
}
