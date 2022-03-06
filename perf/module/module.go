package module

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type PerfModule interface {
	// RegisterCalls registers all possible calls for the module.
	RegisterCalls()

	// Calls random module's RPC method
	CallRandom(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error)
}
