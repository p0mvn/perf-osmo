package epochs

import "github.com/p0mvn/perf-osmo/v2/perf/module"

var _ module.PerfModule                 = (*EpochsModule)(nil)

type EpochsModule struct {

}

func NewPerfModule() module.PerfModule {
	return &EpochsModule{}
}

func (*EpochsModule) RegisterCalls() {

}

func (*EpochsModule) CallRandom() error {
	return nil
}
