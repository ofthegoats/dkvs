////////////////////////////////////////////////////////////
//  This file should contain the smaller methods of Node  //
////////////////////////////////////////////////////////////

package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"

	. "github.com/ofthegoats/dkvs/dkvs-back/rumour"
)

// Update the data present in data.json
func (N *Node) UpdateJson() {
	file, err := json.Marshal(N.Data)
	if err != nil {
		log.Printf("%s: could not write JSON to file\n", fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port))
		return
	}
	err = ioutil.WriteFile("data.json", file, 0644)
	if err != nil {
		log.Printf("%s: could not write JSON to file\n", fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port))
		return
	}
}

func (N *Node) UpdateData(key, newvalue string) {
	oldValue := N.Data[key]
	N.Data[key] = newvalue // Change the existing value, or make a new one
	log.Printf("%s :Updated \"%s\" FROM \"%s\" TO \"%s\"", fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port), key, oldValue, newvalue)
}

func (N *Node) SendNextRound(msg Rumour) {
	if msg.T <= N.MaxRounds {
		msg.T++                                                // increment the round on the message by one
		N.Send(fmt.Sprintf("tcp://127.0.0.1:%d", N.Port), msg) // send to self to continue loop, send via loopback address
		for i := 0; i < N.b; i++ {
			neighbour := N.GetRandomNeighbour()
			log.Printf("%s is sending to %s", fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port), neighbour)
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
		log.Printf("%s: asked to remove neighbour %s which does not exist\n", fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port), neighbour)
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
