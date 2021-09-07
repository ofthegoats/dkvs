package main

import (
	"fmt"
	"log"
	"time"
)

// Request a full state copy from neighbour
func (N *Node) RequestCopy(neighbour string) error {
	requestRumour := Rumour{
		RequestType: FullStateCopyRequest,
		Sender:      fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port),
	}
	err := N.Send(neighbour, requestRumour)
	if err != nil { // failed from this neighbour, not fatal but return error
        log.Printf("%s: Failed to request FSC from %s\n", fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port), neighbour)
		return err
	}
	return nil
}

// infinite! method which performs a RequestCopy every period
func (N *Node) FullStateCopyTimer(period time.Duration) {
	for {
		if len(N.Neighbours) > 0 {
			neighbour := N.GetRandomNeighbour()
			err := N.RequestCopy(neighbour)
			// if there was an error, try at most 5 more times, until there is no error
			for i := 0; i < 5 && err != nil; i++ {
				neighbour = N.GetRandomNeighbour()
				err = N.RequestCopy(neighbour)
			}
		}
		time.Sleep(period)
	}
}
