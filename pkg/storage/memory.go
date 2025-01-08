package storage

import (
	"log"
	"strings"
	"sync"
)

type memoryStorage struct {
	mu           sync.RWMutex
	subscribed   map[string]bool
	transactions map[string][]Transaction
}

// NewMemoryStorage returns an in-memory implementation of Storage.
func NewMemoryStorage() Storage {
	return &memoryStorage{
		subscribed:   make(map[string]bool),
		transactions: make(map[string][]Transaction),
	}
}

func (m *memoryStorage) SubscribeAddress(addr string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	a := strings.ToLower(addr)
	if m.subscribed[a] {
		return false
	}
	m.subscribed[a] = true
	return true
}

func (m *memoryStorage) GetSubscribedAddresses() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var addrs []string
	for k := range m.subscribed {
		addrs = append(addrs, k)
	}
	return addrs
}

func (m *memoryStorage) StoreBlockTransactions(blockNum int, txs []Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, t := range txs {
		from := strings.ToLower(t.From)
		to := strings.ToLower(t.To)
		if m.subscribed[from] || m.subscribed[to] {
			// store in the 'from' address bucket if subscribed
			if m.subscribed[from] {
				m.transactions[from] = append(m.transactions[from], t)
			}

			// store in the 'to' address bucket if subscribed
			if m.subscribed[to] {
				m.transactions[to] = append(m.transactions[to], t)
			}

			log.Println(t)
		}
	}
	return nil
}

func (m *memoryStorage) GetTransactions(addr string) []Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a := strings.ToLower(addr)
	txs := m.transactions[a]

	// copy the slice
	return txs[:]
}
