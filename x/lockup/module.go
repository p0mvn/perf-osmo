package lockup

import "github.com/p0mvn/perf-osmo/v2/perf/module"

var _ module.PerfModule                 = (*LockupModule)(nil)

type LockupModule struct {

}

func NewPerfModule() module.PerfModule {
	return &LockupModule{}
}

func (*LockupModule) RegisterCalls() {

}

func (*LockupModule) CallRandom() error {
	return nil
}
