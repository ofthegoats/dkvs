package rumour

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

	// SuspiciousNode is a value to be used for RequestType
	// It is used to say that a node is suspicious, therefore should be removed from the
	// receiving node's list of neighbours.
	SuspiciousNode = "MARK-NODE-SUSPICIOUS"

	// FullStateCopyRequest is a value to be used for RequestType
	// It is used to ask a node for a full state copy, i.e. copy its data map over one's
	// own. This is used (a) when a node joins the network and (b) randomly for
	// convergence
	FullStateCopyRequest = "FSC-REQUEST"

	// FullStateCopyResponse is a value to be used for RequestType
	// It is used when responding to a node which as requested a full state copy. Requests
	// which use this type must make use the of the FullState value in the Rumour.
	FullStateCopyResponse = "FSC-RESPONSE"

	// GetValueRequest is a value to be used for RequestType
	// It is used to request for the value of a key
	// Rumours using this RequestType need to have a Key value.
	GetValueRequest = "GET-VALUE-REQ"

	// GetValueResponse is a value to be used for RequestType
	// It is used for rumour which do not change the state of a node
	// Instead the rumour is only used to show what the value of a key is.
	GetValueResponse = "GET-VALUE-RESP"
)

// The primary data structure which is communicated between Nodes
type Rumour struct {
	RequestType string // What type of rumour this is, e.g. update data, RTT request ...
	Sender      string // Shows what node to respond to

	Key      string // The key of the piece of data to be updated
	NewValue string // The new value of the piece of data to be updated
	T        int    // The current round this rumour comes from

	RTTTarget   string // If another node is suspicious, fill this with the suspcious socket
	RTTResponse bool   // true if suspcious, else false

	FullState map[string]string // The full state of the data map, to be used for FSCs
}
