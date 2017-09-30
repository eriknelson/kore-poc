package comm

import (
	"fmt"
	"github.com/hegemone/kore-poc/korecomm-go/pkg/config"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

// Engine - Heart of KoreComm, routes ingress and egress messages.
type Engine struct {
	// Messaging buffers
	rawIngressBuffer chan rawIngressBufferMsg
	ingressBuffer    chan ingressBufferMsg
	egressBuffer     chan egressBufferMsg

	// Loaded extensions
	plugins  map[string]*Plugin
	adapters map[string]*Adapter
}

type rawIngressBufferMsg struct {
	AdapterName       string
	RawIngressMessage RawIngressMessage
}

// NOTE: It's possible in the future we'll want some additional metadata
// on this to assist the engine in routing a cmd to a plugin. Right now,
// it's just the cmd less the trigger prefix, which get's matched on the
// plugin's registered manifests
type ingressBufferMsg struct {
	IngressMessage IngressMessage
}

type egressBufferMsg struct {
	Originator    Originator
	EgressMessage EgressMessage
}

// NewEngine - Creates a new work engine.
func NewEngine() *Engine {
	bufferSize := config.GetEngineConfig().BufferSize

	return &Engine{
		rawIngressBuffer: make(chan rawIngressBufferMsg, bufferSize),
		ingressBuffer:    make(chan ingressBufferMsg, bufferSize),
		egressBuffer:     make(chan egressBufferMsg, bufferSize),
		plugins:          make(map[string]*Plugin),
		adapters:         make(map[string]*Adapter),
	}
}

func (e *Engine) LoadExtensions() error {
	log.Info("Loading extensions")
	if err := e.loadPlugins(); err != nil {
		return err
	}
	if err := e.loadAdapters(); err != nil {
		return err
	}
	return nil
}

func (e *Engine) Start() {
	log.Debug("Engine::Start")

	// Spawn listening routines for each adapter
	for _, adapter := range e.adapters {
		adapterCh := make(chan RawIngressMessage)

		go func(adapter *Adapter, adapterch chan RawIngressMessage) {
			// Tell the adapter to start listening and sending messages back via
			// their own ingress channel. Listen should be non-blocking!
			adapter.Listen(adapterCh)

			// Have the listening routine sit and wait for messages back from the
			// adapter. Once received, immediatelly pass them into the raw imsg buffer
			// channel for processing.
			for ribm := range adapterCh {
				e.rawIngressBuffer <- rawIngressBufferMsg{adapter.Name, ribm}
			}
		}(adapter, adapterCh)
	}

	for {
		select {
		case m := <-e.rawIngressBuffer:
			e.handleRawIngress(m)
		case m := <-e.ingressBuffer:
			e.handleIngress(m)
		case m := <-e.egressBuffer:
			e.handleEgress(m)
		}
	}
}

func (e *Engine) handleRawIngress(m rawIngressBufferMsg) {
	go func(ibuff chan<- ingressBufferMsg, m rawIngressBufferMsg) {
		adapterName := m.AdapterName
		rm := m.RawIngressMessage

		if !isCmd(rm.RawContent) {
			return
		}

		if string(rm.RawContent[0]) != CMD_TRIGGER_PREFIX {
			log.Warningf(
				"raw content was flagged as a command, but does not contain trigger prefix, skipping...",
			)
			log.Warning(rm.RawContent)
			return
		}

		content := parseRawContent(rm.RawContent)

		ibuff <- ingressBufferMsg{
			IngressMessage: IngressMessage{
				Originator: Originator{Identity: rm.Identity, AdapterName: adapterName},
				Content:    content,
			},
		}
	}(e.ingressBuffer, m)
}

func parseRawContent(rawContent string) string {
	// NOTE: It's possible in the future we'll want more processing of the raw content
	// some kind of metadata that might be useful for the engine to route the cmd
	// to the plugin. for now
	return rawContent[1:len(rawContent)]
}

func (e *Engine) handleIngress(ibm ingressBufferMsg) {
	im := ibm.IngressMessage
	log.Debugf("Engine::handleIngress: %+v", im)

	go func() {
		cmdMatches := e.applyCmdManifests(im.Content)

		for _, cmdMatch := range cmdMatches {
			delegate := NewCmdDelegate(im, cmdMatch.Submatches)
			cmdMatch.CmdFn(&delegate)
			if delegate.response != "" {
				e.egressBuffer <- egressBufferMsg{
					Originator:    im.Originator,
					EgressMessage: EgressMessage{delegate.response},
				}
			}
		}
	}()
}

type cmdMatch struct {
	CmdFn      CmdFn
	Submatches []string
}

func (e *Engine) applyCmdManifests(content string) []cmdMatch {
	matches := make([]cmdMatch, 0)

	for _, plugin := range e.plugins {
		for _ /*cmdName*/, cmdLink := range plugin.CmdManifest {
			re := cmdLink.Regexp
			subm := re.FindStringSubmatch(content)

			if len(subm) > 0 {
				matches = append(matches, cmdMatch{
					CmdFn:      cmdLink.CmdFn,
					Submatches: subm,
				})
			}
		}
	}

	return matches
}

func (e *Engine) handleEgress(ebm egressBufferMsg) {
	log.Debugf("Engine::handleEgress: %+v", ebm)
	go func() {
		e.adapters[ebm.Originator.AdapterName].SendMessage(ebm.EgressMessage)
	}()
}

// TODO: load{Plugins,Adapters} are almost identical. Should make extension
// loading generic.
func (e *Engine) loadPlugins() error {
	config := config.GetPluginConfig()
	// Check that requested plugins are available in dir, log if not
	log.Infof("Loading plugins from: %v", config.Dir)
	for _, pluginName := range config.Enabled {
		log.Infof("-> %v", pluginName)
		pluginFile := filepath.Join(
			config.Dir,
			fmt.Sprintf("%s.so", pluginName),
		)

		loadedPlugin, err := LoadPlugin(pluginFile)
		if err != nil {
			return err
		}

		e.plugins[loadedPlugin.Name] = loadedPlugin
	}

	log.Info("Successfully loaded plugins:")
	for pluginName, _ := range e.plugins {
		log.Infof("-> %s", pluginName)
	}

	return nil
}

func (e *Engine) loadAdapters() error {
	config := config.GetAdapterConfig()
	// Check that requested adapters are available in dir, log if not
	log.Infof("Loading adapters from: %v", config.Dir)
	for _, adapterName := range config.Enabled {
		log.Infof("-> %v", adapterName)
		adapterFile := filepath.Join(
			config.Dir,
			fmt.Sprintf("%s.so", adapterName),
		)
		log.Infof("file: %s", adapterFile)

		loadedAdapter, err := LoadAdapter(adapterFile)
		if err != nil {
			return err
		}

		loadedAdapter.Init()
		e.adapters[loadedAdapter.Name] = loadedAdapter
	}

	log.Info("Successfully loaded adapters:")
	for adapterName, _ := range e.adapters {
		log.Infof("-> %s", adapterName)
	}
	return nil
}
