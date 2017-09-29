package main

import (
	"fmt"
	"github.com/hegemone/kore-poc/korecomm-go/pkg/comm"
	log "github.com/sirupsen/logrus"
	"time"
)

const intervalSeconds = 1

func Name() string {
	// Reverse DNS identifiers
	return "discord.adapters.kore.nsk.io"
}

func Listen(ch chan<- comm.AdapterIngressMessage) {
	log.Debug("discord.adapters::Listen, spawning ticker")
	go func(ch chan<- comm.AdapterIngressMessage) {
		// In theory, there would be some stop condition here that would done <- true,
		// Ticking perpetually for the sake of demonstration
		ticker := time.NewTicker(intervalSeconds * time.Second)

		count := 1
		tick(count, ch)
		for _ = range ticker.C {
			count++
			tick(count, ch)
		}
	}(ch)
}

func tick(count int, ch chan<- comm.AdapterIngressMessage) {
	log.Debugf("discord.adapters::Tick, sending message -- %d ", count)
	content := fmt.Sprintf("Tick count -> [ %d ]", count)
	ch <- comm.AdapterIngressMessage{Identity: "nsk", Content: content}
}
