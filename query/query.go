package query

import (
	"context"
	"fmt"

	"github.com/p0mvn/perf-osmo/v2/node"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	epochs "github.com/osmosis-labs/osmosis/v7/x/epochs/types"
	lockup "github.com/osmosis-labs/osmosis/v7/x/lockup/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Start() error {
	conn, err := node.NewConnection()
	if err != nil {
		return err
	}

	defer conn.Close()

	reply := &tmservice.GetNodeInfoResponse{}
	if err := conn.Invoke("/cosmos.base.tendermint.v1beta1.Service/GetNodeInfo", reply); err != nil {
		return err
	}
	fmt.Println(reply)

	epochsResp, _, err := conn.InvokeClient(3335437, currEpochRequest)
	if err != nil {
		return err
	}
	fmt.Println(epochsResp)

	lockupResp, _, err := conn.InvokeClient(3335437, getLockupModuleBalance)
	if err != nil {
		return err
	}
	fmt.Println(lockupResp)
	fmt.Println(err)

    return nil
}

func getLockupModuleBalance(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error) {
	lockupClient := lockup.NewQueryClient(grpcConn)

	req := &lockup.ModuleBalanceRequest{}

	return lockupClient.ModuleBalance(ctx, req, grpc.Header(header))
}

func currEpochRequest(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error) {
	epochsClient := epochs.NewQueryClient(grpcConn)

	req := &epochs.QueryCurrentEpochRequest{
		Identifier: "day",
	}
	return epochsClient.CurrentEpoch(ctx, req, grpc.Header(header))
}
