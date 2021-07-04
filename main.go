package main

import "time"

func main() {
	n := NewNode([]string{"tcp://127.0.0.1:9999"}, "tcp://127.0.0.1:9998")
	go n.Gossip()
	time.Sleep(1 * time.Second)
	m := NewNode([]string{"tcp://127.0.0.1:9998"}, "tcp://127.0.0.1:9999")
	testRumour := NewRumour("testKEY", "testVAL")
	m.Send("tcp://127.0.0.1:9998", testRumour)

	// impromptu method to stop from finishing too early
	for {
	}
}
