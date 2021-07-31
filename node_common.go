////////////////////////////////////////////////////////////
//  This file should contain the smaller methods of Node  //
////////////////////////////////////////////////////////////

package main

import (
	"log"
	"math/rand"
)

func (N *Node) UpdateData(key, newvalue string) {
	oldValue := N.Data[key]
	N.Data[key] = newvalue // Change the existing value, or make a new one
	log.Printf("%s :Updated \"%s\" FROM \"%s\" TO \"%s\"", N.socket, key, oldValue, newvalue)
}

func (N *Node) SendNextRound(msg Rumour) {
	if msg.T <= N.MaxRounds {
		msg.T++               // increment the round on the message by one
		N.Send(N.socket, msg) // send to self to continue loop
		for i := 0; i < N.b; i++ {
			neighbour := N.GetRandomNeighbour()
			log.Printf("%s is sending to %s", N.socket, neighbour)
			N.Send(neighbour, msg)
		}
	}
}

func (N *Node) GetRandomNeighbour() string {
	i := rand.Intn(len(N.Neighbours))
	return N.Neighbours[i]
}
