package comm

import (
	log "github.com/sirupsen/logrus"
	"plugin"
)

type Plugin struct {
	Name     string
	Help     string
	Manifest map[string]func(CommandPayload)

	fnName     func() string
	fnHelp     func() string
	fnManifest func() map[string]string
}

type SendResponseFn func(string)

type CommandPayload struct {
	SendResponse   SendResponseFn
	IngressMessage *IngressMessage
	Submatches     []string
}

func LoadPlugin(pluginFile string) (*Plugin, error) {
	p := Plugin{}

	rawPlugin, err := plugin.Open(pluginFile)
	if err != nil {
		return nil, err
	}

	nameSym, err := rawPlugin.Lookup("Name")
	if err != nil {
		return nil, err
	}
	p.fnName = nameSym.(func() string)
	p.Name = p.fnName()

	helpSym, err := rawPlugin.Lookup("Help")
	if err != nil {
		return nil, err
	}
	p.fnHelp = helpSym.(func() string)
	p.Help = p.fnHelp()

	manifestSym, err := rawPlugin.Lookup("Manifest")
	if err != nil {
		return nil, err
	}
	p.fnManifest = manifestSym.(func() map[string]string)

	p.Manifest = make(map[string]func(CommandPayload))
	for cmdRegex, cmdFnName := range p.fnManifest() {
		cmdSym, err := rawPlugin.Lookup(cmdFnName)
		if err != nil {
			log.Error("Error occurred while looking up command for plugin %s: %s", p.Name, err.Error())
			continue
		}

		p.Manifest[cmdRegex] = cmdSym.(func(CommandPayload))
	}

	return &p, nil
}
