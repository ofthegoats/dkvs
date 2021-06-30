package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type Message interface{}

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
// Pass all values sent to Node.Messages channel
func (N *Node) Listen(socket string, messages chan<- Message, wg *sync.WaitGroup) error {
	defer wg.Done()
	listener, err := net.Listen("tcp", N.socket)
	if err != nil { // if listener could not be set up, consider this a fatal error
		log.Fatalln(err)
	}
	defer listener.Close()
	for { // infinite loop accepts all connections
		conn, err := listener.Accept()
		log.Println("connection accepted")
		defer conn.Close()
		if err != nil { // if there is a connection error, log it, but it's not fatal
			log.Println(err)
		}
		go func() { // run this piece of code concurrently:
			for { // accept all messages from the new connection
				msg, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil { // if bufio returns an error, it was the last message
					break
				}
				log.Printf("Listener got the message: %s", msg)
				messages <- msg // send the message into the node's messages channel
			}
		}()
	}
}

func (N *Node) Send(neighbour, message string) error {
	conn, err := net.Dial("tcp", neighbour)
	defer conn.Close()
	if err != nil { // the node might be down, log and return
		log.Println(err)
		return err
	}
	_, err = fmt.Fprintln(conn, message)
	if err != nil { // error on sending to node, unsure of reasons, log and exit
		log.Println(err)
		return err
	}
	return nil // finished with no errors
}

// infinite!
func (N *Node) Gossip() error {
	var wg sync.WaitGroup // wait for all the concurrent procedures to finish before returning
	messages := make(chan Message)
	wg.Add(1)
	go N.Listen(N.socket, messages, &wg)
	for {
		m := <-messages
		log.Printf("got the message: %s\n", m)
	}
}
