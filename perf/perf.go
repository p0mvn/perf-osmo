package perf

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/p0mvn/perf-osmo/v2/perf/module"
	"github.com/p0mvn/perf-osmo/v2/perf/node"
	"github.com/p0mvn/perf-osmo/v2/x/epochs"
	"github.com/p0mvn/perf-osmo/v2/x/lockup"
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
}

func NewManager(host, port string, numConnections, numCallsPerConnection int) *Manager {
	manager := &Manager{
		modules: []module.PerfModule{
			epochs.NewPerfModule(),
			lockup.NewPerfModule(),
		},
		host: host,
		port: port,
		numConnections: numConnections,
		numCallsPerConnection: numCallsPerConnection,
	}

	for _, module := range manager.modules {
		module.RegisterCalls()
	}

	return manager
}

func (m *Manager) Start() error {
	conn, err := node.NewConnection(m.host, m.port)
	if err != nil {
		return err
	}

	defer conn.Close()

	reply := &tmservice.GetNodeInfoResponse{}
	if err := conn.Invoke("/cosmos.base.tendermint.v1beta1.Service/GetNodeInfo", reply); err != nil {
		return err
	}
	fmt.Println(reply)

	epochsResp, _, err := conn.InvokeClient(3335437, epochs.CurrEpochRequest)
	if err != nil {
		return err
	}
	fmt.Println(epochsResp)

	lockupResp, _, err := conn.InvokeClient(3335437, lockup.GetLockupModuleBalance)
	if err != nil {
		return err
	}
	fmt.Println(lockupResp)
	fmt.Println(err)

    return nil
}