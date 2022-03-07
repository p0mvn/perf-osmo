package perf

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/p0mvn/perf-osmo/v2/perf/module"
	"github.com/p0mvn/perf-osmo/v2/perf/node"
	"github.com/p0mvn/perf-osmo/v2/x/epochs"
	"github.com/p0mvn/perf-osmo/v2/x/lockup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Manager struct {
	// All registered modules from which to grab a call to perf test
	modules []module.PerfModule

	// Host
	host string

	// Port
	port string

	// Number of connections
	numConnections int

	// Number of calls to make per connection
	numCallsPerConnection int
	
	// Number of heights below latest to query
	heightsToCover int
}

var _ module.PerfModule                 = (*Manager)(nil)

func NewManager(host, port string, numConnections, numCallsPerConnection, heightsToCover int) *Manager {
	manager := &Manager{
		modules: []module.PerfModule{
			epochs.NewPerfModule(),
			lockup.NewPerfModule(),
		},
		host: host,
		port: port,
		numConnections: numConnections,
		numCallsPerConnection: numCallsPerConnection,
		heightsToCover: heightsToCover,
	}

	manager.RegisterCalls()

	return manager
}

func (m *Manager) RegisterCalls() {
	for _, module := range m.modules {
		module.RegisterCalls()
	}
}

func (m *Manager) CallRandom(grpcConn *grpc.ClientConn, ctx context.Context, header *metadata.MD) (interface{}, error) {
	randN := rand.Intn(len(m.modules))
	return m.modules[randN].CallRandom(grpcConn, ctx, header)
}

func (m *Manager) Start() error {
	rand.Seed(time.Now().UnixNano())
	wg := &sync.WaitGroup{}

	for i :=0; i < m.numConnections; i++ {
		wg.Add(1)
		go m.startConnection(wg)
	}

	wg.Wait()

    return nil
}

func (m *Manager) getLatestHeight() (int64, error) {
	conn, err := node.NewConnection(m.host, m.port)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reply := &tmservice.GetLatestBlockResponse{}
	if err := conn.Invoke("/cosmos.base.tendermint.v1beta1.Service/GetLatestBlock", reply); err != nil {
		return 0, err
	}
	return reply.Block.Header.Height, nil
}

func (m *Manager) startConnection(wg *sync.WaitGroup) error {
	conn, err := node.NewConnection(m.host, m.port)
	if err != nil {
		return err
	}
	defer func ()  {
		conn.Close()
		wg.Done()
	}()

	for i := 0; i < m.numCallsPerConnection; i++ {
		requestHeight, err := m.getLatestHeight()
		if err != nil {
			return err
		}

		heightRand := rand.Intn(100)
		latestHeightChancePercent := 15
		if heightRand > latestHeightChancePercent {
			requestHeight = requestHeight - int64(rand.Intn(m.heightsToCover))
		}

		_, respHeight, err := conn.InvokeClient(int(requestHeight), m.CallRandom)
		if err != nil {
			return err
		}

		if int64(respHeight) != requestHeight {
			return fmt.Errorf("requested height %d, got response height %d", requestHeight, respHeight)
		}
	}
	return nil
}
