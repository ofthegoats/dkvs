package main

import (
	"log"
	"time"
)

// Every `period` amount of time, send an RTT to a random node.
// Follows up the RTT if Send returns an error.
func (N *Node) RTTTimer(period time.Duration, RTTChan chan bool) error {
	RTTRumour := Rumour{
		RequestType: RTTRequest,
		Sender:      N.socket,
	}
	for {
		time.Sleep(period)
		RTTTarget := N.GetRandomNeighbour()
		log.Printf("%s: carrying out RTT on %s\n", N.socket, RTTTarget)
		err := N.Send(RTTTarget, RTTRumour)
		if err != nil {
			log.Printf("%s: node %s did not respond in time\n", N.socket, RTTTarget)
			suspiciousCount := 1 // start with own vote
			for _, neighbour := range N.Neighbours {
				if neighbour == RTTTarget { // we don't check with the failing node
					continue
				}
				err = N.RTTForward(RTTTarget, neighbour)
				if err != nil { // request was not successful, next node
					continue
				}
				result := <-RTTChan
				if result == true {
					suspiciousCount++
				}
			}
			// if a majority of the requested nodes say RTTTarget is suspicious
			if 2*suspiciousCount >= len(N.Neighbours)-1 {
				log.Printf("%s: node %s marked as suspicious\n", N.socket, RTTTarget)
				suspicionRumour := Rumour{
					RequestType: SuspiciousNode,
					RTTTarget:   RTTTarget,
					T:           0,
				}
				// send to self, hack that lets us get this to Gossip() easily, just with localhost
				N.Send(N.socket, suspicionRumour)
                // TODO: tell that node to turn itself off
			}
		}
	}
}

func (N *Node) RTTForward(RTTTarget, neighbour string) error {
	RTTRumour := Rumour{
		RequestType: RTTForward,
		Sender:      N.socket,
		RTTTarget:   RTTTarget,
	}
	err := N.Send(neighbour, RTTRumour)
	return err
}
