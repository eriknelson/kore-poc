package comm

import (
	goplugin "plugin"
)

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

func LoadAdapter(adapterFile string) (*Adapter, error) {
	a := Adapter{}

	rawGoPlugin, err := goplugin.Open(adapterFile)
	if err != nil {
		return nil, err
	}

	nameSym, err := rawGoPlugin.Lookup("Name")
	if err != nil {
		return nil, err
	}
	a.fnName = nameSym.(func() string)
	a.Name = a.fnName()

	listenSym, err := rawGoPlugin.Lookup("Listen")
	if err != nil {
		return nil, err
	}
	a.fnListen = listenSym.(func(chan<- AdapterIngressMessage))

	return &a, nil
}

func processAdapterIngress(adapterName string, am AdapterIngressMessage) (string, IngressMessage) {
	cmd := extractCmd(am.Content)
	return cmd, IngressMessage{
		Originator: Originator{Identity: am.Identity, AdapterName: adapterName},
		Content:    am.Content,
	}
}

func extractCmd(content string) string {
	// TODO: Extract cmd
	return "mock_cmd"
}
