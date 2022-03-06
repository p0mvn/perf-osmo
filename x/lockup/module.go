package lockup

import (
	"context"
	"math/rand"

	"github.com/p0mvn/perf-osmo/v2/perf/module"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ module.PerfModule                 = (*LockupModule)(nil)

type LockupModule struct {
	calls []func(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error)
}

func NewPerfModule() module.PerfModule {
	return &LockupModule{}
}

func (lm *LockupModule) RegisterCalls() {
	lm.calls = append(lm.calls, GetLockupModuleBalance)
}

func (lm *LockupModule) CallRandom(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error) {
	randN := rand.Intn(len(lm.calls))
	return lm.calls[randN](grpcConn, ctx, header)
}
