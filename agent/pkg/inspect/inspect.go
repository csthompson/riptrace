package inspect

import (
	"os"
	"time"

	"github.com/csthompson/riptrace/agent/pkg/types"

	"github.com/google/gops/goprocess"
)

func Inspect() (*types.Profile, error) {
	rtn := types.Profile{}
	rtn.GoProcs = goprocess.FindAll()

	//Retrieve the hostname from the system
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	rtn.Hostname = hostname

	//Retrieve current system time
	rtn.SysTime = time.Now()

	return &rtn, nil
}
