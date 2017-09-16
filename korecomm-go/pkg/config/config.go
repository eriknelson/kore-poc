package config

import (
	//log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

const (
	PluginDirEnvVar  = "KORECOMM_PLUGIN_DIR"
	AdapterDirEnvVar = "KORECOMM_ADAPTER_DIR"
)

type EngineConfig struct {
	BufferSize uint
}

type ExtensionConfig struct {
	Dir     string
	Enabled []string
}

type PluginConfig struct {
	ExtensionConfig
}

type AdapterConfig struct {
	ExtensionConfig
}

// Just fake out the yaml config for now
var mockConfig = map[string]interface{}{
	"engine": EngineConfig{
		BufferSize: 8,
	},
	"plugins": PluginConfig{
		ExtensionConfig: ExtensionConfig{
			Dir: os.Getenv(PluginDirEnvVar),
			Enabled: []string{
				"kore.plugin.bacon",
				"kore.plugin.foo",
			},
		},
	},
	"adapters": AdapterConfig{
		ExtensionConfig: ExtensionConfig{
			Dir: os.Getenv(AdapterDirEnvVar),
			Enabled: []string{
				"kore.adapter.discord",
				"kore.adapter.irc",
			},
		},
	},
}

var _instance *map[string]interface{}
var once sync.Once

func instance() *map[string]interface{} {
	// Threadsafe lazy accessor
	once.Do(func() {
		_instance = loadConfigFile()
	})
	return _instance
}

func GetEngineConfig() EngineConfig {
	return (*instance())["engine"].(EngineConfig)
}

func GetPluginConfig() PluginConfig {
	return (*instance())["plugins"].(PluginConfig)
}

func GetAdapterConfig() AdapterConfig {
	return (*instance())["adapters"].(AdapterConfig)
}

func loadConfigFile() *map[string]interface{} {
	// Load file location from env var, or use default
	// file := GetEnv("KORECOMM_CONFIG") || "/etc/kore/comm_config.yaml"
	return &mockConfig
}
