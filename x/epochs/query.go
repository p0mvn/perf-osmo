package epochs

import (
	"context"

	epochs "github.com/osmosis-labs/osmosis/v7/x/epochs/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CurrEpochRequest(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error) {
	epochsClient := epochs.NewQueryClient(grpcConn)

	req := &epochs.QueryCurrentEpochRequest{
		Identifier: "day",
	}
	return epochsClient.CurrentEpoch(ctx, req, grpc.Header(header))
}
