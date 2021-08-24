package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math"
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
	Timeout    time.Duration     // how long to wait before a TCP request is dropped as a fail
	RTTPeriod  time.Duration     // how long to wait before doing an RTT
	FSCPeriod  time.Duration     // how long to wait before doing a full state copy
	MaxRounds  int

	LANIP string // The local IP for the node
	Port  int    // The port which the node listens on

	key []byte // key used for AES cryptgraphy

	// the number of rounds should be log_b len(Neighbours) + c
	b int // the number of nodes sent to each round
	c int // small configurable value, seeking to make consensus more likely
}

func EncryptAES(key, plaintext []byte) []byte {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce) // fill nonce with cryptographically secure random sequence
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func DecryptAES(key, ciphertext []byte) []byte {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext
}

// Construct a new Node
// knownNeighbours is the neighbours is starts of knowing
// socket is the socket it listens on
// each period is described in the Rumour structure
// b is the number of nodes it should send to each round
// c is the number added to the number of rounds, which should improve probability of consensus
func NewNode(knownNeighbours []string, IP string, key []byte, timeout, rttperiod, fscperiod time.Duration, portNumber, b, c int) Node {
	var n Node
	n.Data = make(map[string]string)
	n.Neighbours = knownNeighbours
	n.Timeout = timeout
	n.RTTPeriod = rttperiod
	n.FSCPeriod = fscperiod
	n.LANIP = IP
	n.Port = portNumber
	n.key = key
	n.b = b
	n.c = c
	n.MaxRounds = int(math.Log(float64(len(n.Neighbours)))/math.Log(float64(n.b))) + n.c // floor (log_b N + c) using base change
	return n
}

// An infinite! procedure that listens on the socket
// All values sent are passed onto the messages channel
func (N *Node) Listen(socket string, messages chan<- Rumour) error {
	lSocket, err := rep.NewSocket() // listenSocket
	if err != nil {                 // If the node fails to establish a socket to listen on, this is fatal
		log.Fatalln(err)
	}
	lSocket.Listen(socket)
	for {
		encrypted_bs, err := lSocket.Recv()
		if err != nil { // failed to read bytes, maybe not be fatal, so don't panic
			log.Println(err)
		}
		bs := DecryptAES(N.key, encrypted_bs)
		// Decode gob structure and send it on messages channel
		reader := bytes.NewReader(bs)
		decoder := gob.NewDecoder(reader)
		var msg Rumour
		err = decoder.Decode(&msg)
		if err != nil { // if decoding failed, decryption must be wrong, implying message was not for this node
			continue
		}
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
	encrypted_buffer := EncryptAES(N.key, buffer.Bytes())
	err = sSocket.Send(encrypted_buffer)
	if err != nil { // sending the bytes failed
		log.Println(err) // not fatal for whole program, but does mean this method failed.
		return err
	}
	return nil
}

// infinite!
func (N *Node) Gossip() error {
	messages := make(chan Rumour)
	listenerSocket := fmt.Sprintf("tcp://%s:%d", N.LANIP, N.Port)
	go N.Listen(listenerSocket, messages)
	RTTChan := make(chan bool)
	go N.RTTTimer(N.RTTPeriod, RTTChan)
	go N.FullStateCopyTimer(N.FSCPeriod)
	for {
		msg := <-messages
		switch msg.RequestType {
		case UpdateData:
			N.UpdateData(msg.Key, msg.NewValue)
			N.SendNextRound(msg)
		case RTTForward:
			err := N.Send(msg.RTTTarget, Rumour{
				RequestType: RTTRequest,
				Sender:      listenerSocket,
			})
			N.Send(msg.Sender, Rumour{
				RequestType: RTTForwardResponse,
				Sender:      listenerSocket,
				RTTResponse: err != nil,
			})
		case RTTForwardResponse:
			RTTChan <- msg.RTTResponse
		case SuspiciousNode:
			N.RemoveNeighbour(msg.RTTTarget)
			N.SendNextRound(msg)
		case FullStateCopyRequest:
			N.Send(msg.Sender, Rumour{
				RequestType: FullStateCopyResponse,
				Sender:      listenerSocket,
				FullState:   N.Data,
			})
		case FullStateCopyResponse:
			N.Data = msg.FullState
		}
	}
}
