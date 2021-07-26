package main

const (
	// UpdateData is a value to be used for RequestType
	// It is used to say that the Rumour asks for a value to be updated for a key
	UpdateData = "UPDATE-DATA"

	// RTTRequest is a value to be used for RequestType.
	// It is used to ask for an RTT from the receiving node.
	RTTRequest = "RTT-REQUEST"

	// RTTForward is a value to be used for RequestType.
	// It is used to ask for the recieving node to send an RTT on a target node.
	RTTForward = "RTT-FORWARD"

	// RTTForwardResponse is a value to be used for RequestType
	// It is used to send the result for an RTTForward Request
	RTTForwardResponse = "RTT-FORWARD-RESPONSE"
)

// The primary data structure which is communicated between Nodes
type Rumour struct {
	RequestType string // What type of rumour this is, e.g. update data, RTT request ...

	Key      string // The key of the piece of data to be updated
	NewValue string // The new value of the piece of data to be updated
	T        int    // The current round this rumour comes from

	RTTTarget string // If another node is suspicious, fill this with the suspcious socket

	RTTResponse bool // true if suspcious, else false

	Sender string // Shows what node to respond to
}
