package main

import (
	"github.com/hegemone/kore-poc/korecomm-go/pkg/comm"
	log "github.com/sirupsen/logrus"
)

var _client *comm.MockPlatformClient

func init() {
	_client = comm.NewMockPlatformClient("discord")
}

func Name() string {
	return "discord.adapters.kore.nsk.io"
}

func Listen(ingressCh chan<- comm.RawIngressMessage) {
	log.Debug("discord.adapters::Listen")

	_client.Connect()

	go func() {
		for clientMsg := range _client.Chat {
			ingressCh <- comm.RawIngressMessage{
				Identity:   clientMsg.User,
				RawContent: clientMsg.Message,
			}
		}
	}()
}

func SendMessage(m string) {
	_client.SendMessage(m)
}
