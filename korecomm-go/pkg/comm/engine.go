package comm

import (
	"fmt"
	"github.com/hegemone/kore-poc/korecomm-go/pkg/config"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	//"github.com/pborman/uuid"
)

// Engine - Heart of KoreComm, routes ingress and egress messages.
type Engine struct {
	ingressBuffer chan IngressMessage
	egressBuffer  chan EgressMessage
	plugins       map[string]*Plugin
	adapters      map[string]*Adapter
}

// NewEngine - Creates a new work engine.
func NewEngine() *Engine {
	bufferSize := config.GetEngineConfig().BufferSize

	return &Engine{
		ingressBuffer: make(chan IngressMessage, bufferSize),
		egressBuffer:  make(chan EgressMessage, bufferSize),
		plugins:       make(map[string]*Plugin),
		adapters:      make(map[string]*Adapter),
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

func (e *Engine) Start() error {
	log.Debug("Engine::Start")

	// TODO: spawn adapter listeners
	// Fan-in to aggregator
	//aggregator := make(chan AdapterIngressMessage)
	// for adapter in adapters
	// adapterchan := (chan<- AdapterIngressMessage)
	// go func(c) {
	//   adapter.Listen(c)
	//   for msg := range c {
	//      aggregator <- msg
	//   }
	// }()
	//select {
	//case msg <- agg:
	//  route to plugin...

	return nil
}

func (e *Engine) SendMessage(originator Originator, responseContent string) {
	log.Debug("Engine::SendMessage -> %s", responseContent)
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

		loadedAdapter, err := LoadAdapter(adapterFile)
		if err != nil {
			return err
		}

		e.adapters[loadedAdapter.Name] = loadedAdapter

		log.Info("Successfully loaded plugins:")
		for adapterName, _ := range e.adapters {
			log.Infof("-> %s", adapterName)
		}
	}
	return nil
}

// Work - is the interface that wraps the basic run method.
//type Work interface {
//Run(token string, msgBuffer chan<- WorkMsg)
//}

// StartNewJob - Starts a job in an new goroutine. returns token, or generated token if an empty token is passed in.
//func (engine *Engine) StartNewJob(token string, work Work) string {
//var jobToken string

//if token == "" {
//jobToken = uuid.New()
//} else {
//jobToken = token
//}
//go work.Run(jobToken, engine.msgBuffer)
//return jobToken
//}

// AttachSubscriber - Attach a subscriber to the engine. Will send the WorkMsg to the subscribers through the message buffer.
//func (engine *Engine) AttachSubscriber(subscriber WorkSubscriber) {
//subscriber.Subscribe(engine.msgBuffer)
//}
