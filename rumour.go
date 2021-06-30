package main

// The dats structure which is communicated between Nodes
type Rumour struct {
	Key      string // The key of the piece of data to be updated
	NewValue string // The new value of the piece of data to be updated
	t        int    // The current round this rumour comes from
}
