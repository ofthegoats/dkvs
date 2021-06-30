package main

import "time"

func main() {
	n := Node{
		socket: "127.0.0.1:9998",
	}
	go n.Gossip()
	time.Sleep(1 * time.Second)
	m := Node{
		socket: "127.0.0.1:9999",
	}
	m.Send("127.0.0.1:9998", "hello world!")

	// impromptu method to stop from finishing too early
	for {
	}
}
