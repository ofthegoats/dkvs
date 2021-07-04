package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"

	"go.nanomsg.org/mangos/v3/protocol/rep"
	"go.nanomsg.org/mangos/v3/protocol/req"

	_ "go.nanomsg.org/mangos/v3/transport/all"
)

// The Node encapsulates the information we are interested in replicating. It also
// communicates with other nodes which it knows of, so as to replicate data with them too.
type Node struct {
	Data       map[string]string // The core part of the Key-Value store, a dictionary.
	Neighbours []string          // List of all known nodes/neighbours.

	socket string // The socket on which this node listens
}

// Construct a new Node
func NewNode(knownNeighbours []string, socket string) Node {
	var n Node
	n.Data = make(map[string]string)
	n.Neighbours = knownNeighbours
	n.socket = socket
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
func (N *Node) Send(neighbour string, message Rumour) error {
	sSocket, err := req.NewSocket() // sendSocket
	if err != nil {                 // failed to establish a socket
		log.Println(err) // not fatal for whole program, but does mean this method failed.
		return err
	}
	sSocket.Dial(neighbour)
	// turn the message into bytes to be sent over the network as gob data
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(message)
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
	for {
		msg := <-messages
		log.Printf("key: %s\nval: %s\n", msg.Key, msg.NewValue)
	}
}
