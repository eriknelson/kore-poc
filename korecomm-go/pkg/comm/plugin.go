package comm

import (
	log "github.com/sirupsen/logrus"
	goplugin "plugin"
	"regexp"
)

type CmdFn func(*CmdDelegate)

type CmdLink struct {
	Regexp *regexp.Regexp
	CmdFn  CmdFn
}

type Plugin struct {
	Name        string
	Help        string
	CmdManifest map[string]CmdLink

	fnName        func() string
	fnHelp        func() string
	fnCmdManifest func() map[string]string
}

func LoadPlugin(pluginFile string) (*Plugin, error) {
	p := Plugin{}

	rawGoPlugin, err := goplugin.Open(pluginFile)
	if err != nil {
		return nil, err
	}

	nameSym, err := rawGoPlugin.Lookup("Name")
	if err != nil {
		return nil, err
	}
	p.fnName = nameSym.(func() string)
	p.Name = p.fnName()

	helpSym, err := rawGoPlugin.Lookup("Help")
	if err != nil {
		return nil, err
	}
	p.fnHelp = helpSym.(func() string)
	p.Help = p.fnHelp()

	cmdManifestSym, err := rawGoPlugin.Lookup("CmdManifest")
	if err != nil {
		return nil, err
	}
	p.fnCmdManifest = cmdManifestSym.(func() map[string]string)

	p.CmdManifest = make(map[string]CmdLink)
	for cmdRegexStr, cmdFnName := range p.fnCmdManifest() {
		cmdFnSym, err := rawGoPlugin.Lookup(cmdFnName)
		if err != nil {
			log.Error("Error occurred while looking up command for plugin %s: %s", p.Name, err.Error())
			continue
		}

		cmdRegex, _ := regexp.Compile(cmdRegexStr)    // TODO: Handle failed regex compilation
		cmdFn := CmdFn(cmdFnSym.(func(*CmdDelegate))) // TODO: Handle failed cast

		// TODO: Error handle more than one command named the same thing
		p.CmdManifest[cmdFnName] = CmdLink{
			Regexp: cmdRegex,
			CmdFn:  cmdFn,
		}
	}

	return &p, nil
}
