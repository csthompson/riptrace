package types

import "github.com/go-delve/delve/service/api"

type TraceInfo struct {
	api.DebuggerState
}

func (self TraceInfo) GetLocals() []api.Variable {
	return self.CurrentThread.BreakpointInfo.Locals
}

func (self TraceInfo) GetVariables() []api.Variable {
	return self.CurrentThread.BreakpointInfo.Variables
}
