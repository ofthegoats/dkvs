package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/rep"
	"go.nanomsg.org/mangos/v3/protocol/req"

	_ "go.nanomsg.org/mangos/v3/transport/all"
)

// The Node encapsulates the information we are interested in replicating. It also
// communicates with other nodes which it knows of, so as to replicate data with them too.
type Node struct {
	Data       map[string]string // The core part of the Key-Value store, a dictionary.
	Neighbours []string          // List of all known nodes/neighbours.
	Timeout    time.Duration
	RTTPeriod  time.Duration

	socket string // The socket on which this node listens

	// the number of rounds should be log_b len(Neighbours) + c
	b int // the number of nodes sent to each round
	c int // small configurable value, seeking to make consensus more likely
}

// Construct a new Node
// knownNeighbours is the neighbours is starts of knowing
// socket is the socket it listens on
// b is the number of nodes it should send to each round
// c is the number added to the number of rounds, which should improve probability of consensus
func NewNode(knownNeighbours []string, socket string, timeout time.Duration, period time.Duration, b, c int) Node {
	var n Node
	n.Data = make(map[string]string)
	n.Neighbours = knownNeighbours
	n.Timeout = timeout
	n.RTTPeriod = period
	n.socket = socket
	n.b = b
	n.c = c
	return n
}

// An infinite! procedure that listens on the socket
// All values sent are passed onto the messages channel
func (N *Node) Listen(socket string, messages chan<- Rumour, wg *sync.WaitGroup) error {
	defer wg.Done()
	lSocket, err := rep.NewSocket() // listenSocket
	if err != nil {                 // If the node fails to establish a socket to listen on, this is fatal
		log.Fatalln(err)
	}
	lSocket.Listen(socket)
	for {
		bs, err := lSocket.Recv()
		if err != nil { // failed to read bytes, maybe not be fatal, so don't panic
			log.Println(err)
		}
		// Decode gob structure and send it on messages channel
		reader := bytes.NewReader(bs)
		decoder := gob.NewDecoder(reader)
		var msg Rumour
		decoder.Decode(&msg)
		messages <- msg
	}
}

// Given a neighbour and a rumour, send the rumour to that neighbour.
// Does not select a random neigbour or cycle.
func (N *Node) Send(neighbour string, rumour Rumour) error {
	sSocket, err := req.NewSocket() // sendSocket
	if err != nil {                 // failed to establish a socket
		log.Println(err) // not fatal for whole program, but does mean this method failed.
		return err
	}
	sSocket.SetOption(mangos.OptionSendDeadline, N.Timeout)
	sSocket.Dial(neighbour)
	// turn the message into bytes to be sent over the network as gob data
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(rumour)
	if err != nil { // failed to encode, do not send
		log.Println(err)
		return err
	}
	err = sSocket.Send(buffer.Bytes())
	if err != nil { // sending the bytes failed
		log.Println(err) // not fatal for whole program, but does mean this method failed.
		return err
	}
	return nil
}

// infinite!
func (N *Node) Gossip() error {
	var wg sync.WaitGroup // wait for all the concurrent procedures to finish before returning
	messages := make(chan Rumour)
	wg.Add(1)
	go N.Listen(N.socket, messages, &wg)
	wg.Add(1)
	go N.RTTTimer(N.RTTPeriod)
	maxRounds := int(math.Log(float64(len(N.Neighbours)))/math.Log(float64(N.b))) + N.c // floor (log_b N + c) using base change
	log.Printf("maxRounds = %d", maxRounds)
	for {
		msg := <-messages
		switch msg.RequestType {
		case UpdateData:
			// Update the value stored in the node
			oldValue := N.Data[msg.Key]
			N.Data[msg.Key] = msg.NewValue // Change the existing value, or make a new one
			log.Printf("%s :Updated \"%s\" FROM \"%s\" TO \"%s\"", N.socket, msg.Key, oldValue, msg.NewValue)
			if msg.T <= maxRounds {
				msg.T++ // increment the round on the message by one
				for i := 0; i < N.b; i++ {
					neighbour := N.GetRandomNeighbour()
					log.Printf("%s is sending to %s", N.socket, neighbour)
					N.Send(neighbour, msg)
				}
			}
		}
	}
}

func (N *Node) GetRandomNeighbour() string {
	i := rand.Intn(len(N.Neighbours))
	return N.Neighbours[i]
}
