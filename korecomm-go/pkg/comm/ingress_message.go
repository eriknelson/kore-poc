package comm

type Originator struct {
	Identity    string
	AdapterName string
}

type IngressMessage struct {
	Content    string
	Originator Originator
}

// RawIngressMessage - Raw, unprocessed message passed from the adapter to the
// engine. Has not yet been parsed to determine if the message is a cmd or not.
type RawIngressMessage struct {
	Identity   string
	RawContent string
}
