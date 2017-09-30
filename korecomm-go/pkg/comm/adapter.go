package comm

import (
	"fmt"
	goplugin "plugin"
	"regexp"
)

var (
	CMD_TRIGGER_PREFIX = "!"
	CMD_REGEX, _       = regexp.Compile(fmt.Sprintf("^%s\\S*($| )", CMD_TRIGGER_PREFIX))
)

type Adapter struct {
	Name string

	//fnInit        func(MessageReceivedCallback)
	fnSendMessage func(string)
	fnName        func() string
	fnListen      func(chan<- RawIngressMessage)
}

func (a *Adapter) Listen(inChan chan<- RawIngressMessage) {
	// Possibly some common logic an Adapter might want to do instead of having
	// the engine call the raw plugin listen directly
	// NOTE: Engine has already handled spawning the listen routines in their
	// own goroutines, so they're running concurrently and/or in parallel.
	// TODO: Engine probably needs to handle the case of adapters being poorly
	// written and immediately having their channels close or leaked. fnListen
	// is expected to be long lived.
	a.fnListen(inChan)
}

func (a *Adapter) SendMessage(emsg EgressMessage) {
	// TODO: General processing
	a.fnSendMessage(emsg.Serialize())
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
	a.fnListen = listenSym.(func(chan<- RawIngressMessage))

	sendMessageSym, err := rawGoPlugin.Lookup("SendMessage")
	if err != nil {
		return nil, err
	}
	a.fnSendMessage = sendMessageSym.(func(string))

	return &a, nil
}

func isCmd(rawContent string) bool {
	return CMD_REGEX.MatchString(rawContent)
}

func extractCmd(content string) string {
	// TODO: Extract cmd
	return "mock_cmd"
}
