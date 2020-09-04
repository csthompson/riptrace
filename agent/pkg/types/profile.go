package types

import (
	"time"

	"github.com/google/gops/goprocess"
)

type Profile struct {
	GoProcs  []goprocess.P
	Hostname string
	SysTime  time.Time
}
