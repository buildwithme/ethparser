package storage

// Transaction captures minimal TX data.
type Transaction struct {
	Hash        string
	From        string
	To          string
	BlockNumber int
	Value       string
}

// Storage is an interface for storing and retrieving TXs.
type Storage interface {
	// SubscribeAddress adds an address for tracking.
	SubscribeAddress(addr string) bool

	// GetSubscribedAddresses returns all subscribed addresses.
	GetSubscribedAddresses() []string

	// StoreBlockTransactions does an atomic insertion of all TXs for a block.
	// Only stores if TX's 'from' or 'to' is subscribed.
	StoreBlockTransactions(blockNum int, txs []Transaction) error

	// GetTransactions returns stored TXs for a specific address, sorted by block.
	GetTransactions(addr string) []Transaction
}
