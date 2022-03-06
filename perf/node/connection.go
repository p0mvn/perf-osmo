package node

import (
	"context"
	"fmt"
	"strconv"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Connection struct {
	grpcConn *grpc.ClientConn
}

func NewConnection(host, port string) (*Connection, error) {
	grpcConn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port),grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Connection{
		grpcConn: grpcConn,
	}, nil
}

func (c *Connection) InvokeClient(height int, makeRequestCb func (*grpc.ClientConn,  context.Context, *metadata.MD) (interface{}, error)) (interface{}, int, error) {
	var header metadata.MD
	ctxWithHeight := metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, strconv.Itoa(height))

	resp, err := makeRequestCb(c.grpcConn, ctxWithHeight, &header)
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

func (c *Connection) Invoke(method string, reply interface{}) error {
	if err := c.grpcConn.Invoke(context.Background(), method, nil, reply, grpc.EmptyCallOption{}); err != nil {
		return err
	}
	return nil
}

func (c *Connection) Close() error {
	return c.grpcConn.Close()
}

