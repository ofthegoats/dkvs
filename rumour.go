package main

// The primary data structure which is communicated between Nodes
type Rumour struct {
	Key      string // The key of the piece of data to be updated
	NewValue string // The new value of the piece of data to be updated
	T        int    // The current round this rumour comes from
}

// Makes a new rumour, to be sent around the network
func NewRumour(key string, newval string) Rumour {
	r := Rumour{
		Key:      key,
		NewValue: newval,
		T:        0, // upon creation, 0 rounds have been passed.
	}
	return r
}
