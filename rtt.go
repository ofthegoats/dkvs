package main

import "time"

// Every `period` amount of time, send an RTT to a random node.
// Follows up the RTT if the SendRTT returns true
func (N *Node) RTTTimer(period time.Duration) error {
	RTTRumour := Rumour{
		RequestType: RTTRequest,
	}
	for {
		time.Sleep(period)
		RTTTarget := N.GetRandomNeighbour()
		err := N.Send(RTTTarget, RTTRumour)
		if err != nil {
			// TODO: send a request to some neighbours, see if they also can't reach the node.
		}
	}
}
