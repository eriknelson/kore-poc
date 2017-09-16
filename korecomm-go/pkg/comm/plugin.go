package comm

import (
	"plugin"
	//log "github.com/sirupsen/logrus"
)

type Plugin struct {
	Hello func(IngressMessage)
}

func LoadPlugin(pluginFile string) (*Plugin, error) {
	p := Plugin{}

	rawPlugin, err := plugin.Open(pluginFile)
	if err != nil {
		return nil, err
	}

	helloSym, err := rawPlugin.Lookup("Hello")
	if err != nil {
		return nil, err
	}
	p.Hello = helloSym.(func(IngressMessage))

	return &p, nil
}
