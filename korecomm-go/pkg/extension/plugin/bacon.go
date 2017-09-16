package main

import (
	kc "github.com/hegemone/kore-poc/korecomm-go/pkg/comm"

	log "github.com/sirupsen/logrus"
)

var pluginName = "kore.plugin.bacon"

func Hello(msg kc.IngressMessage) {
	log.Infof("[%s]: Hello, this is %v, got content: %v", pluginName, pluginName, msg.Content)
}
