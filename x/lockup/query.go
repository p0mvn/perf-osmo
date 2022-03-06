package lockup

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	lockup "github.com/osmosis-labs/osmosis/v7/x/lockup/types"
)


func GetLockupModuleBalance(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error) {
	lockupClient := lockup.NewQueryClient(grpcConn)

	req := &lockup.ModuleBalanceRequest{}

	return lockupClient.ModuleBalance(ctx, req, grpc.Header(header))
}
