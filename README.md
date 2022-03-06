# perf-osmo

This is a simple application to performance test an Osmosis node.

Allows adding rpc endpoints to query by modules (similar to Cosmos SDK).

Sidenote: could easily be adapted to any Cosmos SDK node in the future.

To use:

- Place `config.yml` into `$HOME/.perf-osmo`:
```
---
# See perf.go for interpretation
perf:
  host: <node ip>
  port: 9090
  numConnections: 1
  numCallsPerConnection: 1
  heightsToCover: 1000

```

- `go install`

- `perf-osmos start`
