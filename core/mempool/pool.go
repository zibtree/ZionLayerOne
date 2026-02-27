package mempool

import (
	"errors"
	"sort"
	"sync"

	"github.com/zionlayer/zionlayer/core/transaction"
)

const (
	MaxPoolSize = 10_000
)

var (
	ErrPoolFull    = errors.New("mempool is full")
	ErrDuplicateTx = errors.New("duplicate transaction")
)

// Pool is a thread-safe transaction pool.
type Pool struct {
	mu  sync.RWMutex
	txs map[[32]byte]*transaction.Tx
}

// NewPool creates an empty mempool.
func NewPool() *Pool {
	return &Pool{
		txs: make(map[[32]byte]*transaction.Tx),
	}
}

// Add inserts a transaction into the pool.
func (p *Pool) Add(tx *transaction.Tx) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.txs) >= MaxPoolSize {
		return ErrPoolFull
	}
	h := tx.Hash()
	if _, exists := p.txs[h]; exists {
		return ErrDuplicateTx
	}
	p.txs[h] = tx
	return nil
}

// Pop removes and returns up to n transactions, sorted by gas price descending.
func (p *Pool) Pop(n int) []*transaction.Tx {
	p.mu.Lock()
	defer p.mu.Unlock()

	all := make([]*transaction.Tx, 0, len(p.txs))
	for _, tx := range p.txs {
		all = append(all, tx)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].GasPrice.Cmp(all[j].GasPrice) > 0
	})
	if n > len(all) {
		n = len(all)
	}
	selected := all[:n]
	for _, tx := range selected {
		delete(p.txs, tx.Hash())
	}
	return selected
}

// Size returns the number of pending transactions.
func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.txs)
}
