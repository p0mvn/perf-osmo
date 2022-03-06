package module

type PerfModule interface {
	// RegisterCalls registers all possible calls for the module.
	RegisterCalls()

	// Calls random module's RPC method
	CallRandom() error
}
