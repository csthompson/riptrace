package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/csthompson/riptrace/agent/pkg/inspect"
	"github.com/csthompson/riptrace/agent/service/natsvc"
)

type ProfileHandler struct {
	NatsClient *natsvc.NatsClient
}

func (self *ProfileHandler) Get(subj string, reply string, m interface{}) {
	log.Info("Handling get profile request")
	profile, err := inspect.Inspect()
	if err != nil {
		log.Error("Error retrieving profile ", err)
	}
	self.NatsClient.Publish(reply, profile)
}
