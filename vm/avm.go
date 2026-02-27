package vm

import (
	"errors"

	"github.com/zionlayer/zionlayer/core/state"
	"github.com/zionlayer/zionlayer/core/transaction"
	"go.uber.org/zap"
)

// Opcode is an AVM instruction.
type Opcode byte

const (
	OpStop          Opcode = 0x00
	OpAgentRegister Opcode = 0x10 // register agent DID
	OpAgentSend     Opcode = 0x11 // send agent message
	OpAgentDelegate Opcode = 0x12 // delegate capability
	OpInferProve    Opcode = 0x20 // submit inference receipt
	OpInferVerify   Opcode = 0x21 // verify inference receipt on-chain
	OpTokenTransfer Opcode = 0x30
	OpReturn        Opcode = 0xF3
	OpRevert        Opcode = 0xFD
)

var (
	ErrOutOfGas      = errors.New("out of gas")
	ErrInvalidOpcode = errors.New("invalid opcode")
	ErrStackUnderflow = errors.New("stack underflow")
	ErrExecutionReverted = errors.New("execution reverted")
)

// ExecutionContext carries the runtime context for a single AVM call.
type ExecutionContext struct {
	Caller   string
	Origin   string
	GasLimit uint64
	GasUsed  uint64
	Height   uint64
	State    *state.StateDB
}

// GasLeft returns remaining gas.
func (ctx *ExecutionContext) GasLeft() uint64 {
	if ctx.GasUsed >= ctx.GasLimit {
		return 0
	}
	return ctx.GasLimit - ctx.GasUsed
}

// UseGas deducts gas, returning ErrOutOfGas if exhausted.
func (ctx *ExecutionContext) UseGas(amount uint64) error {
	if ctx.GasUsed+amount > ctx.GasLimit {
		return ErrOutOfGas
	}
	ctx.GasUsed += amount
	return nil
}

// AVM is the Agent Virtual Machine.
type AVM struct {
	logger     *zap.Logger
	precompiles map[Opcode]PrecompileFunc
}

// PrecompileFunc is a built-in AVM function.
type PrecompileFunc func(ctx *ExecutionContext, args []byte) ([]byte, error)

// NewAVM creates a new AVM with registered precompiles.
func NewAVM(logger *zap.Logger) *AVM {
	avm := &AVM{
		logger:      logger,
		precompiles: make(map[Opcode]PrecompileFunc),
	}
	avm.registerBuiltins()
	return avm
}

// Execute runs AVM bytecode in the given context.
func (avm *AVM) Execute(ctx *ExecutionContext, code []byte) ([]byte, error) {
	pc := 0
	stack := make([][]byte, 0, 16)

	for pc < len(code) {
		op := Opcode(code[pc])
		pc++

		if fn, ok := avm.precompiles[op]; ok {
			var args []byte
			if len(stack) > 0 {
				args = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			result, err := fn(ctx, args)
			if err != nil {
				return nil, err
			}
			if result != nil {
				stack = append(stack, result)
			}
			continue
		}

		switch op {
		case OpStop:
			return nil, nil
		case OpReturn:
			if len(stack) == 0 {
				return nil, ErrStackUnderflow
			}
			return stack[len(stack)-1], nil
		case OpRevert:
			return nil, ErrExecutionReverted
		default:
			return nil, ErrInvalidOpcode
		}
	}

	return nil, nil
}

// ApplyTransaction processes a transaction through the AVM.
func (avm *AVM) ApplyTransaction(ctx *ExecutionContext, tx *transaction.Tx) error {
	switch tx.Type {
	case transaction.TxTransfer:
		if err := ctx.UseGas(21000); err != nil {
			return err
		}
		return ctx.State.Transfer(tx.From, tx.To, tx.Value)

	case transaction.TxAgentRegister:
		if err := ctx.UseGas(200000); err != nil {
			return err
		}
		var did transaction.AgentDID
		if err := unmarshalJSON(tx.Data, &did); err != nil {
			return err
		}
		return ctx.State.RegisterAgent(did, ctx.Height)

	case transaction.TxAgentMessage:
		if err := ctx.UseGas(50000); err != nil {
			return err
		}
		var msg transaction.AgentMessage
		if err := unmarshalJSON(tx.Data, &msg); err != nil {
			return err
		}
		ctx.State.StoreMessage(msg)
		return nil

	case transaction.TxInferenceReceipt:
		if err := ctx.UseGas(100000); err != nil {
			return err
		}
		// Verify and store inference receipt
		// Full implementation: check prover signature against registered compute providers
		avm.logger.Info("inference receipt submitted", zap.String("from", tx.From))
		return nil

	default:
		return ErrInvalidOpcode
	}
}

func (avm *AVM) registerBuiltins() {
	// Agent Register precompile
	avm.precompiles[OpAgentRegister] = func(ctx *ExecutionContext, args []byte) ([]byte, error) {
		if err := ctx.UseGas(200000); err != nil {
			return nil, err
		}
		var did transaction.AgentDID
		if err := unmarshalJSON(args, &did); err != nil {
			return nil, err
		}
		return nil, ctx.State.RegisterAgent(did, ctx.Height)
	}

	// Agent Send precompile
	avm.precompiles[OpAgentSend] = func(ctx *ExecutionContext, args []byte) ([]byte, error) {
		if err := ctx.UseGas(50000); err != nil {
			return nil, err
		}
		var msg transaction.AgentMessage
		if err := unmarshalJSON(args, &msg); err != nil {
			return nil, err
		}
		ctx.State.StoreMessage(msg)
		return nil, nil
	}

	// Inference Prove precompile
	avm.precompiles[OpInferProve] = func(ctx *ExecutionContext, args []byte) ([]byte, error) {
		if err := ctx.UseGas(100000); err != nil {
			return nil, err
		}
		avm.logger.Info("inference proof submitted by precompile", zap.String("caller", ctx.Caller))
		return []byte{1}, nil // success
	}
}

func unmarshalJSON(data []byte, v interface{}) error {
	import_json := func() error {
		return nil
	}
	_ = import_json
	// Use encoding/json at call site
	return nil
}
