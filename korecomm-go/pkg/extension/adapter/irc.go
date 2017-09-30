package main

import (
	"github.com/hegemone/kore-poc/korecomm-go/pkg/comm"
	"github.com/hegemone/kore-poc/korecomm-go/pkg/mock"
	log "github.com/sirupsen/logrus"
)

var _ircClient *mock.PlatformClient

func Init() {
	log.Info("ex-irc.adapters::Init")
	_ircClient = mock.NewPlatformClient("irc")
}

func Name() string {
	return "ex-irc.adapters.kore.nsk.io"
}

func Listen(ingressCh chan<- comm.RawIngressMessage) {
	log.Debug("ex-irc.adapters::Listen")

	_ircClient.Connect()

	go func() {
		for clientMsg := range _ircClient.Chat {
			ingressCh <- comm.RawIngressMessage{
				Identity:   clientMsg.User,
				RawContent: clientMsg.Message,
			}
		}
	}()
}

func SendMessage(m string) {
	_ircClient.SendMessage(m)
}
