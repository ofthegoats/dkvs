package main

import (
	"log"
	"time"
)

func main() {
	n := NewNode([]string{"tcp://127.0.0.1:9999"}, "tcp://127.0.0.1:9998", 1*time.Second, 2, 1)
	go n.Gossip()
	time.Sleep(1 * time.Second)
	m := NewNode([]string{"tcp://127.0.0.1:9998"}, "tcp://127.0.0.1:9999", 1*time.Second, 2, 1)
	go m.Gossip()
	time.Sleep(1 * time.Second)

	r1 := NewRumour("key1", "val1")
	n.Send("tcp://127.0.0.1:9998", r1) // socket DOES exist
	n.Send("tcp://127.0.0.1:9997", r1) // socket does not exist
	n.Send("tcp://127.0.0.1:9996", r1) // socket does not exist

	time.Sleep(1 * time.Second) // give time to catch up

	log.Printf("%v\n", n.Data)
	log.Printf("%v\n", m.Data)
}
