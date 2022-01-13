package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/ofthegoats/dkvs/dkvs-back/node"
	. "github.com/ofthegoats/dkvs/dkvs-back/rumour"
)

// syntax: dkvs-front $PORT $KEY $COMMAND [ $KEY | $KEY $VALUE | $SOCKET ]

func checkArgs() error {
	err := errors.New("not enough arguments")
	length := len(os.Args)
	if length < 3 { // always necessary: port, key, command
		return err
	}
	switch {
	case os.Args[3] == "add-neighbour" && length < 4:
		return err
	case os.Args[3] == "remove-neighbour" && length < 4:
		return err
	case os.Args[3] == "get-value" && length < 4:
		return err
	case os.Args[3] == "set-value" && length < 5:
		return err
	}
	return nil
}

func main() {
	// check if the correct amount of arguments has been passed
	if err := checkArgs(); err != nil {
		log.Fatalln(err)
	}
	port := os.Args[1] // port on which the *backend* is listening
	key := []byte(os.Args[2])
	command := os.Args[3]
	// fairly arbitrary values taken from back/main.go
	// TODO: replace with unix sockets
	tcptimeout := 1 * time.Second
	rttperiod := 300 * time.Second
	fscperiod := 300 * time.Second
	n := NewNode(
		[]string{}, "localhost",
		key,
		tcptimeout, rttperiod, fscperiod,
		9875, 5, 2,
	)
	messages := make(chan Rumour)
	listener := fmt.Sprintf("tcp://localhost:%d", 9875)
	backend := fmt.Sprintf("tcp://localhost:%s", port)
	go n.Listen(listener, messages)
	switch command {
	case "list-neighbours":
		err := n.Send(backend, Rumour{
			RequestType: GetNeighboursRequest,
			Sender:      listener,
		})
		if err != nil {
			log.Println(err)
		} else {
			neighbours := <-messages
			fmt.Printf("%v\n", neighbours.Neighbours)
		}
	case "add-neighbour":
		newNeighbour := os.Args[4]
		err := n.Send(backend, Rumour{
			RequestType: AddNeighbourRequest,
			Sender:      listener,
			NewValue:    newNeighbour,
		})
		if err != nil {
			log.Println(err)
		} else {
			time.Sleep(1 * time.Second) // let send catch up
		}
	case "remove-neighbour":
		toRemove := os.Args[4]
		err := n.Send(backend, Rumour{
			RequestType: DeleteNeighbourRequest,
			Sender:      listener,
			NewValue:    toRemove,
		})
		if err != nil {
			log.Println(err)
		} else {
			time.Sleep(1 * time.Second) // let send catch up
		}
	case "list-values":
		err := n.RequestCopy(backend)
		if err != nil {
			log.Println(err)
		} else {
			fsc := <-messages
			fmt.Printf("%v\n", fsc.FullState)
		}
	case "get-value":
		k := os.Args[4]
		err := n.Send(backend, Rumour{
			RequestType: GetValueRequest,
			Sender:      listener,
			Key:         k,
		})
		if err != nil {
			log.Println(err)
		} else {
			v := <-messages
			fmt.Printf("%v\n", v.NewValue)
		}
	case "set-value":
		k := os.Args[4]
		v := os.Args[5]
		err := n.Send(backend, Rumour{
			RequestType: UpdateData,
			Sender:      listener,
			Key:         k,
			NewValue:    v,
			T:           0,
		})
		if err != nil {
			log.Println(err)
		} else {
			time.Sleep(1 * time.Second) // let send catch up
		}
	case "die":
		err := n.Send(backend, Rumour{
			RequestType: DieRequest,
			Sender:      listener,
		})
		if err != nil {
			log.Println(err)
		} else {
			time.Sleep(1 * time.Second) // let send catch up
		}
	default:
		log.Fatalf("not a valid command: %s\n", command)
	}
}
