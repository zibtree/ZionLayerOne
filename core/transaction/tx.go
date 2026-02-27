package transaction

import (
	"crypto/sha256"
	"encoding/json"
	"math/big"
)

// TxType classifies the transaction.
type TxType uint8

const (
	TxTransfer          TxType = iota // standard token transfer
	TxAgentRegister                   // register a new agent DID
	TxAgentMessage                    // agent-to-agent message
	TxAgentDelegate                   // delegate capability to another agent
	TxDeployContract                  // deploy AVM contract
	TxCallContract                    // call AVM contract
	TxInferenceReceipt                // submit verifiable inference proof
	TxValidatorStake                  // stake tokens as validator
	TxValidatorUnstake                // unstake tokens
)

// Capability represents a named agent capability.
type Capability struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// AgentDID is a decentralized identifier anchored on-chain.
type AgentDID struct {
	ID           string            `json:"id"`            // did:agc:0x...
	Controller   string            `json:"controller"`    // owner address (hex)
	Capabilities []Capability      `json:"capabilities"`
	PublicKey    []byte            `json:"publicKey"`
	Metadata     map[string]string `json:"metadata"`
}

// MessageType classifies agent messages.
type MessageType string

const (
	MsgTask     MessageType = "TASK"
	MsgResult   MessageType = "RESULT"
	MsgDelegate MessageType = "DELEGATE"
	MsgRevoke   MessageType = "REVOKE"
)

// AgentMessage is an on-chain structured message between agents.
type AgentMessage struct {
	From    string      `json:"from"`
	To      string      `json:"to"`
	Type    MessageType `json:"type"`
	Payload []byte      `json:"payload"`
	Nonce   uint64      `json:"nonce"`
}

// InferenceReceipt is a verifiable proof of AI inference.
type InferenceReceipt struct {
	AgentID    string `json:"agentId"`
	ModelHash  []byte `json:"modelHash"`  // IPFS CID bytes
	InputHash  []byte `json:"inputHash"`
	OutputHash []byte `json:"outputHash"`
	Timestamp  int64  `json:"timestamp"`
	ProverSig  []byte `json:"proverSig"`
}

// Tx is a signed transaction on ZionLayer.
type Tx struct {
	Type      TxType          `json:"type"`
	From      string          `json:"from"`    // sender address (hex)
	To        string          `json:"to"`      // recipient address (hex)
	Value     *big.Int        `json:"value"`   // $ZIO in smallest unit
	Gas       uint64          `json:"gas"`
	GasPrice  *big.Int        `json:"gasPrice"`
	Nonce     uint64          `json:"nonce"`
	Data      json.RawMessage `json:"data"`    // type-specific payload
	Signature []byte          `json:"sig"`
}

// Hash returns the SHA-256 hash of the transaction (excluding signature).
func (tx *Tx) Hash() [32]byte {
	cp := *tx
	cp.Signature = nil
	data, _ := json.Marshal(cp)
	return sha256.Sum256(data)
}

// NewTransferTx creates a basic token transfer transaction.
func NewTransferTx(from, to string, value *big.Int, nonce uint64, gasPrice *big.Int) *Tx {
	return &Tx{
		Type:     TxTransfer,
		From:     from,
		To:       to,
		Value:    value,
		Gas:      21000,
		GasPrice: gasPrice,
		Nonce:    nonce,
	}
}

// NewAgentRegisterTx creates an agent registration transaction.
func NewAgentRegisterTx(from string, did AgentDID, nonce uint64, gasPrice *big.Int) *Tx {
	data, _ := json.Marshal(did)
	return &Tx{
		Type:     TxAgentRegister,
		From:     from,
		Gas:      200000,
		GasPrice: gasPrice,
		Nonce:    nonce,
		Data:     data,
	}
}

// NewAgentMessageTx creates an agent message transaction.
func NewAgentMessageTx(from string, msg AgentMessage, nonce uint64, gasPrice *big.Int) *Tx {
	data, _ := json.Marshal(msg)
	return &Tx{
		Type:     TxAgentMessage,
		From:     from,
		Gas:      50000,
		GasPrice: gasPrice,
		Nonce:    nonce,
		Data:     data,
	}
}

// NewInferenceReceiptTx creates an inference receipt submission transaction.
func NewInferenceReceiptTx(from string, receipt InferenceReceipt, nonce uint64, gasPrice *big.Int) *Tx {
	data, _ := json.Marshal(receipt)
	return &Tx{
		Type:     TxInferenceReceipt,
		From:     from,
		Gas:      100000,
		GasPrice: gasPrice,
		Nonce:    nonce,
		Data:     data,
	}
}
