package block

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/zionlayer/zionlayer/core/transaction"
)

// Header contains the block metadata.
type Header struct {
	Version        uint32
	Height         uint64
	Timestamp      int64
	PrevHash       [32]byte
	StateRoot      [32]byte
	TxRoot         [32]byte
	AgentRoot      [32]byte // merkle root of agent state trie
	ValidatorAddr  []byte
	Signature      []byte
}

// Block is a full block including header and transactions.
type Block struct {
	Header Header
	Txs    []*transaction.Tx
}

// NewBlock creates a new block with the given header fields.
func NewBlock(height uint64, prevHash [32]byte, validatorAddr []byte, txs []*transaction.Tx) *Block {
	return &Block{
		Header: Header{
			Version:       1,
			Height:        height,
			Timestamp:     time.Now().UnixNano(),
			PrevHash:      prevHash,
			ValidatorAddr: validatorAddr,
		},
		Txs: txs,
	}
}

// Hash returns the SHA-256 hash of the block header.
func (b *Block) Hash() [32]byte {
	data, _ := json.Marshal(b.Header)
	return sha256.Sum256(data)
}

// GenesisBlock creates the genesis block.
func GenesisBlock() *Block {
	return &Block{
		Header: Header{
			Version:   1,
			Height:    0,
			Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano(),
		},
		Txs: []*transaction.Tx{},
	}
}
