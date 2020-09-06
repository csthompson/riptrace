package publishers

import (
	"github.com/csthompson/riptrace/agent/pkg/types"
	"github.com/csthompson/riptrace/agent/service/debugger"
	"github.com/csthompson/riptrace/agent/service/natsvc"

	log "github.com/sirupsen/logrus"
)

type TracePublisher struct {
	NatsClient *natsvc.NatsClient
	Debugger   *debugger.Debugger
	Terminate  bool
}

func (self *TracePublisher) Publish() {
	log.Info("Trace publisher started")
	go func() {
		for {
			log.Info("Waiting...")
			st, err := self.Debugger.Continue()
			if err != nil {
				log.Error(err)
			}
			t := types.TraceInfo{*st}
			//Terminate before publishing
			if self.Terminate {
				break
			}
			self.NatsClient.Publish("debugger.trace", t)
		}
	}()
}

func (self *TracePublisher) Shutdown() {
	log.Info("Trace publisher shutting down")
	self.Terminate = true
}
