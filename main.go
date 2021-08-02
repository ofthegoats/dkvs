package main

import (
	"log"
	"time"
)

func main() {
	n := NewNode([]string{"tcp://127.0.0.1:9999", "tcp://localhost:9996"}, "tcp://127.0.0.1:9998", 1*time.Second, 2*time.Second, 2, 10)
	go n.Gossip()
	time.Sleep(1 * time.Second)

    m := NewNode([]string{"tcp://127.0.0.1:9998", "tcp://localhost:9996"}, "tcp://127.0.0.1:9999", 1*time.Second, 2*time.Second, 2, 10)
	go m.Gossip()
	time.Sleep(1 * time.Second)

	log.Printf("%s: maxrounds = %d\n", n.socket, n.MaxRounds)
	log.Printf("%s: maxrounds = %d\n", m.socket, m.MaxRounds)

	r1 := Rumour{
		RequestType: UpdateData,
		Key:         "key1",
		NewValue:    "val1",
		T:           0,
	}
	n.Send("tcp://127.0.0.1:9998", r1) // socket DOES exist
	n.Send("tcp://127.0.0.1:9997", r1) // socket does NOT exist
	n.Send("tcp://127.0.0.1:9996", r1) // socket does NOT exist

	time.Sleep(1 * time.Second) // give time to catch up

	for {
		log.Printf("%s: %v, %v\n", n.socket, n.Data, n.Neighbours)
		log.Printf("%s: %v, %v\n", m.socket, m.Data, n.Neighbours)
		time.Sleep(5 * time.Second)
	}
}
