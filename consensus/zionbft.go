package consensus

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/zionlayer/zionlayer/core/block"
	"github.com/zionlayer/zionlayer/core/state"
	"github.com/zionlayer/zionlayer/core/transaction"
	"go.uber.org/zap"
)

const (
	BlockTime       = 2 * time.Second
	MinValidatorStake = 10_000 // in ZIO base units (×10^18)
	BlockReward     = 5       // ZIO per block
)

var (
	ErrInvalidBlock     = errors.New("invalid block")
	ErrInvalidSignature = errors.New("invalid block signature")
	ErrUnknownValidator = errors.New("unknown validator")
)

// Validator represents a staked network validator.
type Validator struct {
	Address    string
	PublicKey  []byte
	Stake      *big.Int
	PoIScore   float64 // Proof-of-Intelligence score
	VotingPower int64
}

// ZionBFT is the hybrid PoS + PoI consensus engine.
type ZionBFT struct {
	mu         sync.RWMutex
	validators map[string]*Validator
	state      *state.StateDB
	logger     *zap.Logger
	height     uint64
	tip        *block.Block

	// channels
	blockCh chan *block.Block
	quitCh  chan struct{}
}

// NewZionBFT creates a new consensus engine.
func NewZionBFT(stateDB *state.StateDB, logger *zap.Logger) *ZionBFT {
	return &ZionBFT{
		validators: make(map[string]*Validator),
		state:      stateDB,
		logger:     logger,
		blockCh:    make(chan *block.Block, 64),
		quitCh:     make(chan struct{}),
	}
}

// AddValidator registers a validator with the consensus engine.
func (e *ZionBFT) AddValidator(v *Validator) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	minStake := new(big.Int).Mul(big.NewInt(MinValidatorStake), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	if v.Stake.Cmp(minStake) < 0 {
		return errors.New("stake below minimum")
	}
	e.validators[v.Address] = v
	e.logger.Info("validator registered", zap.String("addr", v.Address))
	return nil
}

// Start begins block production.
func (e *ZionBFT) Start(proposerAddr string, txPool <-chan []*transaction.Tx) {
	go e.runProposer(proposerAddr, txPool)
}

// Stop halts the consensus engine.
func (e *ZionBFT) Stop() {
	close(e.quitCh)
}

// Blocks returns the channel of finalized blocks.
func (e *ZionBFT) Blocks() <-chan *block.Block {
	return e.blockCh
}

// ValidateBlock checks block validity.
func (e *ZionBFT) ValidateBlock(b *block.Block) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	v, ok := e.validators[string(b.Header.ValidatorAddr)]
	if !ok {
		return ErrUnknownValidator
	}
	_ = v // signature verification would go here

	if b.Header.Height != e.height+1 {
		return ErrInvalidBlock
	}
	prevHash := e.tip.Hash()
	if b.Header.PrevHash != prevHash {
		return ErrInvalidBlock
	}
	return nil
}

// runProposer produces blocks at BlockTime intervals.
func (e *ZionBFT) runProposer(addr string, txPool <-chan []*transaction.Tx) {
	ticker := time.NewTicker(BlockTime)
	defer ticker.Stop()

	for {
		select {
		case <-e.quitCh:
			return
		case <-ticker.C:
			var txs []*transaction.Tx
			select {
			case batch := <-txPool:
				txs = batch
			default:
				txs = []*transaction.Tx{}
			}

			e.mu.Lock()
			var prevHash [32]byte
			if e.tip != nil {
				prevHash = e.tip.Hash()
			}
			b := block.NewBlock(e.height+1, prevHash, []byte(addr), txs)
			// In production: compute state root, sign block, broadcast for votes
			e.height++
			e.tip = b
			e.mu.Unlock()

			e.applyBlockReward(addr)
			e.logger.Info("block proposed", zap.Uint64("height", b.Header.Height), zap.Int("txs", len(txs)))

			select {
			case e.blockCh <- b:
			default:
			}
		}
	}
}

func (e *ZionBFT) applyBlockReward(validatorAddr string) {
	reward := new(big.Int).Mul(
		big.NewInt(BlockReward),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
	)
	acc := e.state.GetAccount(validatorAddr)
	newBal := new(big.Int).Add(acc.Balance, reward)
	e.state.SetBalance(validatorAddr, newBal)
}

// VotingPower computes a validator's voting power from stake + PoI score.
func (e *ZionBFT) VotingPower(v *Validator) int64 {
	stakeScore := new(big.Int).Div(v.Stake, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)).Int64()
	poiBoost := int64(v.PoIScore * 100)
	return stakeScore + poiBoost
}
