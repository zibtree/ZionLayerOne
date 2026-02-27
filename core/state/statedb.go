package state

import (
	"encoding/json"
	"errors"
	"math/big"
	"sync"

	"github.com/zionlayer/zionlayer/core/transaction"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrAgentNotFound = errors.New("agent not found")
	ErrAgentAlreadyRegistered = errors.New("agent already registered")
)

// Account holds the state of an address.
type Account struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
	Nonce   uint64   `json:"nonce"`
	Code    []byte   `json:"code,omitempty"` // AVM bytecode if contract
}

// AgentRecord stores on-chain agent metadata.
type AgentRecord struct {
	DID          transaction.AgentDID `json:"did"`
	RegisteredAt uint64               `json:"registeredAt"` // block height
	MessageCount uint64               `json:"messageCount"`
	Active       bool                 `json:"active"`
}

// StateDB is the in-memory world state.
// In production this wraps an iavl MerkleTrie.
type StateDB struct {
	mu       sync.RWMutex
	accounts map[string]*Account
	agents   map[string]*AgentRecord // keyed by DID.ID
	messages []transaction.AgentMessage
}

// NewStateDB initializes a fresh StateDB.
func NewStateDB() *StateDB {
	return &StateDB{
		accounts: make(map[string]*Account),
		agents:   make(map[string]*AgentRecord),
	}
}

// GetAccount returns the account for an address, creating it if needed.
func (s *StateDB) GetAccount(addr string) *Account {
	s.mu.RLock()
	defer s.mu.RUnlock()
	acc, ok := s.accounts[addr]
	if !ok {
		return &Account{Address: addr, Balance: big.NewInt(0)}
	}
	return acc
}

// SetBalance sets the balance for an address.
func (s *StateDB) SetBalance(addr string, balance *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	acc := s.getOrCreate(addr)
	acc.Balance = new(big.Int).Set(balance)
}

// Transfer moves value from one address to another.
func (s *StateDB) Transfer(from, to string, value *big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	src := s.getOrCreate(from)
	if src.Balance.Cmp(value) < 0 {
		return ErrInsufficientBalance
	}
	dst := s.getOrCreate(to)
	src.Balance.Sub(src.Balance, value)
	dst.Balance.Add(dst.Balance, value)
	return nil
}

// RegisterAgent registers a new AgentDID on-chain.
func (s *StateDB) RegisterAgent(did transaction.AgentDID, blockHeight uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.agents[did.ID]; exists {
		return ErrAgentAlreadyRegistered
	}
	s.agents[did.ID] = &AgentRecord{
		DID:          did,
		RegisteredAt: blockHeight,
		Active:       true,
	}
	return nil
}

// GetAgent returns the agent record for a DID.
func (s *StateDB) GetAgent(didID string) (*AgentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.agents[didID]
	if !ok {
		return nil, ErrAgentNotFound
	}
	return rec, nil
}

// StoreMessage appends an agent message to the log.
func (s *StateDB) StoreMessage(msg transaction.AgentMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
	if rec, ok := s.agents[msg.From]; ok {
		rec.MessageCount++
	}
}

// Snapshot serializes the full state to JSON (simplified; production uses MerkleTrie).
func (s *StateDB) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	type snap struct {
		Accounts map[string]*Account      `json:"accounts"`
		Agents   map[string]*AgentRecord  `json:"agents"`
	}
	return json.Marshal(snap{Accounts: s.accounts, Agents: s.agents})
}

func (s *StateDB) getOrCreate(addr string) *Account {
	if acc, ok := s.accounts[addr]; ok {
		return acc
	}
	acc := &Account{Address: addr, Balance: big.NewInt(0)}
	s.accounts[addr] = acc
	return acc
}
