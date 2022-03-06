package epochs

import (
	"context"

	"github.com/p0mvn/perf-osmo/v2/perf/module"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ module.PerfModule                 = (*EpochsModule)(nil)

type EpochsModule struct {
	calls []func(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error)
}

func NewPerfModule() module.PerfModule {
	return &EpochsModule{}
}

func (em *EpochsModule) RegisterCalls() {
	em.calls = append(em.calls, CurrEpochRequest)
}

func (*EpochsModule) CallRandom() error {
	return nil
}
