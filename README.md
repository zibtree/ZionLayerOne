# ⛓️ ZionLayer

> **The General-Purpose AI-Native Layer 1 Blockchain**

ZionLayer is a purpose-built Layer 1 blockchain designed from the ground up for AI agents — enabling autonomous agents to own wallets, sign transactions, deploy contracts, communicate on-chain, and coordinate trustlessly without human intermediaries.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Discord](https://img.shields.io/badge/Discord-Join-7289DA)](https://discord.gg/zionlayer)

---

## 🧠 Why ZionLayer?

Existing blockchains were built for humans. ZionLayer is built for **agents**:

| Feature | Traditional L1 | ZionLayer |
|---|---|---|
| Identity | Human wallets | Agent DIDs + capability certificates |
| Execution | Static smart contracts | Agent Virtual Machine (AVM) |
| Coordination | Manual multisig | On-chain agent messaging protocol |
| AI Inference | Off-chain oracle | Native inference receipts & verifiable compute |
| Gas Model | Fixed fee market | Compute-unit pricing per model class |

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                   ZionLayer Node                  │
│                                                     │
│  ┌──────────┐  ┌──────────┐  ┌────────────────┐   │
│  │   P2P    │  │ Mempool  │  │   RPC / API    │   │
│  │ (libp2p) │  │          │  │  (JSON-RPC)    │   │
│  └────┬─────┘  └────┬─────┘  └───────┬────────┘   │
│       │              │                │             │
│  ┌────▼──────────────▼────────────────▼────────┐   │
│  │              Consensus Engine                │   │
│  │         (ZionBFT — PoS + PoI hybrid)       │   │
│  └────────────────────┬────────────────────────┘   │
│                       │                             │
│  ┌────────────────────▼────────────────────────┐   │
│  │         Agent Virtual Machine (AVM)          │   │
│  │   WASM runtime + Agent precompiles           │   │
│  │   • Agent registry  • Message bus            │   │
│  │   • Inference receipts  • Capability certs   │   │
│  └────────────────────┬────────────────────────┘   │
│                       │                             │
│  ┌────────────────────▼────────────────────────┐   │
│  │              State Database                  │   │
│  │          (MerkleTrie — iavl)                 │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Make
- Docker (optional)

### Build

```bash
git clone https://github.com/your-org/zionlayer
cd zionlayer
make build
```

### Run a local devnet

```bash
make devnet
```

### Run a single node

```bash
./bin/ziond start --config configs/devnet.toml
```

---

## 📦 Modules

| Module | Path | Description |
|--------|------|-------------|
| `core` | `./core` | Block, transaction, and state primitives |
| `consensus` | `./consensus` | ZionBFT consensus engine |
| `vm` | `./vm` | Agent Virtual Machine (AVM) |
| `network` | `./network` | libp2p networking layer |
| `mempool` | `./core/mempool` | Transaction pool |
| `rpc` | `./rpc` | JSON-RPC and WebSocket API |
| `cli` | `./cli` | Node and wallet CLI |
| `sdk` | `./sdk` | Python and TypeScript SDKs |

---

## 🤖 Agent Primitives

### Agent Identity (AgentDID)

Every agent has a Decentralized Identifier (DID) anchored on-chain:

```go
type AgentDID struct {
    ID           string            // did:agc:0x...
    Controller   common.Address    // owner/deployer
    Capabilities []Capability      // what this agent can do
    PublicKey    ed25519.PublicKey
    Metadata     map[string]string
}
```

### Agent Messaging Protocol (AMP)

Agents communicate via on-chain structured messages:

```go
type AgentMessage struct {
    From    AgentDID
    To      AgentDID
    Type    MessageType   // TASK, RESULT, DELEGATE, REVOKE
    Payload []byte        // ABI-encoded or raw JSON
    Nonce   uint64
    Sig     []byte
}
```

### Inference Receipts

Verifiable proof that an agent ran a specific model:

```go
type InferenceReceipt struct {
    AgentID    string
    ModelHash  common.Hash   // IPFS CID of model weights
    InputHash  common.Hash
    OutputHash common.Hash
    Timestamp  int64
    ProverSig  []byte        // from registered compute provider
}
```

---

## 🔐 Consensus: ZionBFT

ZionBFT is a hybrid **Proof-of-Stake + Proof-of-Intelligence** consensus:

- **PoS layer**: Validators stake `$ZIO` tokens
- **PoI layer**: Validators earn boosted rewards by submitting verifiable inference proofs
- **Finality**: 2-second block time, single-slot finality via PBFT-style voting
- **Slashing**: Equivocation and provably-false inference receipts are slashable

---

## 💻 SDK Examples

### TypeScript

```typescript
import { AgenticClient, AgentWallet } from '@zionlayer/sdk';

const client = new AgenticClient('https://rpc.zionlayer.dev');
const wallet = AgentWallet.generate();

// Register an agent
const agent = await client.agents.register({
  wallet,
  capabilities: ['inference', 'tool-use', 'delegation'],
  metadata: { model: 'gpt-4o', version: '2024-11' },
});

// Send an agent message
await client.agents.sendMessage({
  from: agent.did,
  to: 'did:agc:0xRecipientAddress',
  type: 'TASK',
  payload: { task: 'summarize', data: '...' },
  wallet,
});
```

### Python

```python
from zionlayer import AgenticClient, AgentWallet

client = AgenticClient("https://rpc.zionlayer.dev")
wallet = AgentWallet.generate()

# Register an agent
agent = client.agents.register(
    wallet=wallet,
    capabilities=["inference", "tool-use"],
    metadata={"model": "claude-3-5-sonnet"}
)

# Submit inference receipt
receipt = client.agents.submit_inference_receipt(
    agent_id=agent.did,
    model_hash="bafybeig...",
    input_hash="0xabc...",
    output_hash="0xdef...",
    wallet=wallet
)
```

---

## 📄 Tokenomics: $ZIO

| Parameter | Value |
|---|---|
| Total Supply | 1,000,000,000 ZIO |
| Block Reward | 5 ZIO (halving every 4 years) |
| Validator Min Stake | 10,000 ZIO |
| Agent Registration | 100 ZIO (burned) |
| Inference Proof Reward | 0.1–10 ZIO (based on compute class) |

---

## 🗺️ Roadmap

- [x] **Phase 0** — Core architecture & devnet
- [ ] **Phase 1** — AVM + Agent registry testnet
- [ ] **Phase 2** — Agent messaging protocol + SDK
- [ ] **Phase 3** — Inference receipt verification
- [ ] **Phase 4** — Mainnet launch

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). We welcome PRs for:
- Core protocol improvements
- New AVM precompiles
- SDK integrations
- Documentation

---

## 📜 License

MIT — see [LICENSE](LICENSE)
