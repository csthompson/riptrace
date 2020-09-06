package main

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/csthompson/riptrace/agent/internal/handlers"
	"github.com/csthompson/riptrace/agent/pkg/types"
	"github.com/csthompson/riptrace/agent/service/debugger"
	"github.com/csthompson/riptrace/agent/service/natsvc"
)

var agentProcess sync.WaitGroup

//The agent can only have one active debugger at a time (agent -> program)
var Debugger debugger.Debugger

func main() {

	osSigCh := make(chan os.Signal)
	signal.Notify(osSigCh, os.Interrupt, syscall.SIGTERM)

	agentProcess.Add(1)

	natsClient := natsvc.GetClient("localhost:4222")

	//Defer and release resources gracefully
	defer func() {
		natsClient.Client.Close()
		agentProcess.Done()
	}()

	//Handle requests to the hosts profile
	profileHandler := handlers.ProfileHandler{NatsClient: &natsClient}
	natsClient.RegisterTopicHandler("profile.get", profileHandler.Get)

	//Handle requests to the debugger
	debugHandler := &handlers.DebuggerHandler{NatsClient: &natsClient}
	natsClient.RegisterTopicHandler("debugger.attach", debugHandler.Attach)
	natsClient.RegisterTopicHandler("debugger.createBreakpoint", debugHandler.CreateBreakpoint)

	//handle sigterm
	go func() {
		select {
		case <-osSigCh:
			log.Error("Shutting down riptrace agent ", debugHandler.Debugger)
			if debugHandler.Debugger != nil {
				debugHandler.Shutdown()
			}
			natsClient.Client.Close()
			os.Exit(1)
		}
	}()

	//THIS IS ALL TEST DATA. NOT PART OF SERVER
	result := types.Profile{}
	natsClient.Request("profile.get", nil, &result)

	pid := 0
	for _, p := range result.GoProcs {
		if strings.Contains(p.Path, "popquote/main") {
			pid = p.PID
		}
		log.Println(p.PID, ' ', p.Exec, ' ', p.Path)
	}

	var attach bool
	log.Println(pid)
	natsClient.Request("debugger.attach", pid, &attach)

	log.Info("Attach status ", attach)

	bp := types.Breakpoint{
		File: "/home/csthompson/workspace/go/src/github.com/csthompson/popquote/internal/handlers/annotations.go",
		Line: 70,
	}
	natsClient.Request("debugger.createBreakpoint", bp, &attach)

	natsClient.Client.Subscribe("debugger.trace", func(trace *types.TraceInfo) {
		if trace != nil {
			for _, v := range trace.GetLocals() {
				log.Info(v.Name, v.Type)
			}
		}
	})

	agentProcess.Wait()
}
