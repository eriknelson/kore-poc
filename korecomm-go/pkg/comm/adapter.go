package comm

type Adapter struct {
	Name string

	//fnInit        func(MessageReceivedCallback)
	//fnSendMessage func(EgressMessage)
	fnName   func() string
	fnListen func(chan<- AdapterIngressMessage)
}

func (a *Adapter) Listen(inChan chan<- AdapterIngressMessage) {
	// Possibly some common logic an Adapter might want to do instead of having
	// the engine call the raw plugin listen directly
	// NOTE: Engine has already handled spawning the listen routines in their
	// own goroutines, so they're running concurrently and/or in parallel.
	// TODO: Engine probably needs to handle the case of adapters being poorly
	// written and immediately having their channels close or leaked. fnListen
	// is expected to be long lived.
	a.fnListen(inChan)
}
