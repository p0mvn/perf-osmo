package query

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	epochs "github.com/osmosis-labs/osmosis/v7/x/epochs/types"
	lockup "github.com/osmosis-labs/osmosis/v7/x/lockup/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func Start() error {

	// Create a connection to the gRPC server.
	grpcTendermintConn, err := grpc.Dial("104.248.92.191:9090",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer grpcTendermintConn.Close()
	

	reply := &tmservice.GetNodeInfoResponse{}
	if err := grpcTendermintConn.Invoke(context.Background(), "/cosmos.base.tendermint.v1beta1.Service/GetNodeInfo", nil, reply, grpc.EmptyCallOption{}); err != nil {
		return err
	}
	fmt.Println(reply)

    // // Create a connection to the gRPC server.
    // grpcConn, err := grpc.Dial("104.248.92.191:9090",grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	return err
	// }
    // defer grpcConn.Close()

	// epochsResp, _, err := makeRequest(grpcConn, 3335437, currEpochRequest)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(epochsResp)

	// lockupResp, _, err := makeRequest(grpcConn, 3335437, getLockupModuleBalance)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(lockupResp)
	// fmt.Println(err)

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

func makeRequest(grpcConn *grpc.ClientConn, height int, makeRequestCb func (*grpc.ClientConn,  context.Context, *metadata.MD) (interface{}, error)) (interface{}, int, error) {
	var header metadata.MD
	ctxWithHeight := metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, strconv.Itoa(height))

	resp, err := makeRequestCb(grpcConn, ctxWithHeight, &header)
	if err != nil {
		return nil, 0, err
	}

	blockHeightUnparsed := header.Get(grpctypes.GRPCBlockHeightHeader)
	if len(blockHeightUnparsed) > 1 {
		return nil, 0, fmt.Errorf("invalid block height header: %s", blockHeightUnparsed)
	}

	blockHeight, err := strconv.Atoi(blockHeightUnparsed[0])
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing block height header: %s", blockHeightUnparsed)
	}

	return resp, blockHeight, nil
}
