package comm

type EgressSender interface {
	SendEgress(EgressMessage)
}

type EgressMessage struct {
	Content string
}

func (e *EgressMessage) Serialize() string {
	// NOTE: Might need to expand on this in the future
	return e.Content
}
