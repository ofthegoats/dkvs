////////////////////////////////////////////////////////////
//  This file should contain the smaller methods of Node  //
////////////////////////////////////////////////////////////

package main

import (
	"log"
	"math"
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

func (N *Node) RemoveNeighbour(neighbour string) {
	index := 0
	found := false
	for i, n := range N.Neighbours {
		if n == neighbour {
			index = i
			found = true
			break
		}
	}
	if found == false { // the neigbour to be removed does not exist
		log.Printf("%s: asked to remove neighbour %s which does not exist\n", N.socket, neighbour)
	} else {
		// swap the neighbour to be removed with the first, then remove the first
		// we can do this because order does not matter
		N.Neighbours[0], N.Neighbours[index] = N.Neighbours[index], N.Neighbours[0]
		N.Neighbours = N.Neighbours[1:]
		N.MaxRounds = int(math.Log(float64(len(N.Neighbours)))/math.Log(float64(N.b))) + N.c // need to update MaxRounds after changing neighbours
	}
}

func (N *Node) GetRandomNeighbour() string {
	i := rand.Intn(len(N.Neighbours))
	return N.Neighbours[i]
}
