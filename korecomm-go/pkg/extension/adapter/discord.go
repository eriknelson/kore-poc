package main

import (
	"github.com/hegemone/kore-poc/korecomm-go/pkg/comm"
	"github.com/hegemone/kore-poc/korecomm-go/pkg/mock"
	log "github.com/sirupsen/logrus"
)

var _discordClient *mock.PlatformClient

func Init() {
	log.Info("ex-discord.adapters::Init")
	_discordClient = mock.NewPlatformClient("discord")
}

func Name() string {
	return "ex-discord.adapters.kore.nsk.io"
}

func Listen(ingressCh chan<- comm.RawIngressMessage) {
	log.Debug("ex-discord.adapters::Listen")

	_discordClient.Connect()

	go func() {
		for clientMsg := range _discordClient.Chat {
			ingressCh <- comm.RawIngressMessage{
				Identity:   clientMsg.User,
				RawContent: clientMsg.Message,
			}
		}
	}()
}

func SendMessage(m string) {
	_discordClient.SendMessage(m)
}
