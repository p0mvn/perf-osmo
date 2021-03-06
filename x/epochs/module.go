package epochs

import (
	"context"
	"math/rand"

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

func (em *EpochsModule) CallRandom(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error) {
	randN := rand.Intn(len(em.calls))
	return em.calls[randN](grpcConn, ctx, header)
}
