package handlers

import (
	"github.com/csthompson/riptrace/agent/internal/publishers"
	"github.com/csthompson/riptrace/agent/pkg/types"
	"github.com/csthompson/riptrace/agent/service/debugger"
	"github.com/csthompson/riptrace/agent/service/natsvc"
	log "github.com/sirupsen/logrus"
)

type DebuggerHandler struct {
	NatsClient *natsvc.NatsClient
	Debugger   *debugger.Debugger
	Publishers map[string]types.Publisher
}

const (
	TRACE_PUBLISHER = "TRACE_PUBLISHER"
)

//Send a signal to all active publishers to shutdown
// Detatch the debugger
func (self *DebuggerHandler) Shutdown() {
	for _, p := range self.Publishers {
		p.Shutdown()
	}
	self.Debugger.Shutdown()
}

func (self *DebuggerHandler) Attach(subj string, reply string, msg int) {
	log.Info("Handling debugger attach request")

	//First check to see if we have an active debugger
	if self.Debugger != nil {
		log.Info("Active riptrace process already running")
		return
	}

	pid := msg

	d, err := debugger.New(pid)
	if err != nil {
		log.Error("Error initiating riptrace on PID ", pid, " : ", err)
		return
	}

	log.Info("Successfully Attached to PID: ", d.Delve.ProcessPid())

	self.Debugger = d

	if err != nil {
		self.NatsClient.Publish(reply, false)
		return
	} else {
		self.NatsClient.Publish(reply, true)
		return
	}
}

func (self *DebuggerHandler) CreateBreakpoint(subj string, reply string, bp *types.Breakpoint) {

	err := self.Debugger.CreateBreakpoint(bp.File, bp.Line)

	//We added a breakpoint, so now we have a trace publisher
	if _, ok := self.Publishers[TRACE_PUBLISHER]; !ok {
		if self.Publishers == nil {
			self.Publishers = map[string]types.Publisher{}
		}
		p := publishers.TracePublisher{
			NatsClient: self.NatsClient,
			Debugger:   self.Debugger,
		}
		self.Publishers[TRACE_PUBLISHER] = &p
		self.Publishers[TRACE_PUBLISHER].Publish()
	}

	if err != nil {
		self.NatsClient.Publish(reply, false)
		return
	} else {
		self.NatsClient.Publish(reply, true)
		return
	}
}
