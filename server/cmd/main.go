package main

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/csthompson/riptrace/agent/pkg/types"
	"github.com/csthompson/riptrace/agent/service/natsvc"
)

var serverProcess sync.WaitGroup

func main() {

	serverProcess.Add(1)

	natsClient := natsvc.GetClient("localhost:4222")

	//THIS IS ALL TEST DATA. NOT PART OF THE OVERALL SERVER DESIGN.
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
		File: "<TEST FILE>",
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

	serverProcess.Wait()
}
