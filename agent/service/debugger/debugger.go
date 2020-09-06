package debugger

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
)

type Debugger struct {
	Delve *debugger.Debugger
}

//Attach to an existing process
func New(pid int) (*Debugger, error) {
	config := &debugger.Config{
		AttachPid:   pid,
		Backend:     "default",
		WorkingDir:  ".",
		ExecuteKind: debugger.ExecutingOther,
	}
	d, err := debugger.New(config, []string{})
	if err != nil {
		return nil, err
	}

	debugger := Debugger{
		Delve: d,
	}
	return &debugger, nil
}

//Cleanup and let program continue executing
func (self *Debugger) Shutdown() {
	log.Println("Removing breakpoints")
	self.Delve.Command(&api.DebuggerCommand{Name: api.Halt})
	for _, bp := range self.Delve.Breakpoints() {
		log.Info("Removing breakpoint ", bp.Addrs)
		_, err := self.Delve.ClearBreakpoint(bp)
		if err != nil {
			log.Error("Could not clear breakpoint ", err)
		}
	}
	self.Delve.Detach(false)
}

func (self *Debugger) Continue() (*api.DebuggerState, error) {
	st, err := self.Delve.Command(&api.DebuggerCommand{Name: api.Continue})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return st, nil
}

func (self *Debugger) GetState() {
	st, err := self.Delve.State(false)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(st.CurrentThread)
}

func (self *Debugger) CreateBreakpoint(file string, line int) error {
	bp := &api.Breakpoint{
		File:       file,
		Line:       line,
		Tracepoint: true,
		Stacktrace: 1,
		LoadLocals: &api.LoadConfig{
			FollowPointers:     true,
			MaxVariableRecurse: 1,
			MaxStructFields:    -1,
		},
	}
	_, err := self.Delve.CreateBreakpoint(bp)
	if err != nil {
		return err
	}
	return nil
}
