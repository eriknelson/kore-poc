package main

import (
	"fmt"
	"github.com/hegemone/kore-poc/korecomm-go/pkg/comm"
	log "github.com/sirupsen/logrus"
)

func Name() string {
	// Reverse DNS identifiers
	return "bacon.plugins.kore.nsk.io"
}

func Manifest() map[string]string {
	return map[string]string{
		`/bacon$/`:       "CommandBacon",
		`/bacon\s+(\S+)`: "CommandBaconGift",
	}
}

func Help() string {
	return "Usage: !bacon [user]"
}

func CommandBacon(p comm.CommandPayload) {
	log.Info("bacon::CommandBacon")

	msg := p.IngressMessage
	identity := msg.Originator.Identity

	response := fmt.Sprintf(
		"gives %s a strip of delicious bacon.", identity,
	)

	p.SendResponse(response)
}

func CommandBaconGift(p comm.CommandPayload) {
	log.Info("bacon::CommandBaconGift")

	msg := p.IngressMessage
	identity := msg.Originator.Identity
	toUser := p.Submatches[1]

	response := fmt.Sprintf(
		"gives %s a strip of delicious bacon as a gift from %v",
		toUser, identity,
	)

	p.SendResponse(response)
}
