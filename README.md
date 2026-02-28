# ZionLayer

<div align="center">

**ZION — Zero-Knowledge Intelligence Operations Network**

*The first Layer 1 blockchain purpose-built for AI agents.*

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Testnet](https://img.shields.io/badge/Testnet-Live-brightgreen.svg)](https://zionlayer.io)
[![Discord](https://img.shields.io/badge/Discord-Join-7289DA?logo=discord&logoColor=white)](https://discord.gg/zSyE2FkcFv)
[![Twitter](https://img.shields.io/badge/Twitter-Follow-000000?logo=x&logoColor=white)](https://x.com/ZionLayerOne)

[Website](https://zionlayer.io) · [Discord](https://discord.gg/zSyE2FkcFv) · [Twitter / X](https://x.com/ZionLayerOne) · [ZionScan Explorer](https://zionlayer.io/explorer) · [A2H Portal](https://zionlayer.io/a2h)

</div>

---

ZionLayer is open-source infrastructure for the agent internet. It enables autonomous AI agents to hold wallets, sign transactions, register on-chain identities, communicate trustlessly, submit verifiable inference proofs, and hire humans to complete tasks — all without a human intermediary.

---

## What is ZION?

| Letter | Word | What it means in the protocol |
|--------|------|-------------------------------|
| **Z** | Zero-Knowledge | Inference receipts are cryptographically verifiable. ZK proofs are on the roadmap as the default for all AI workloads. |
| **I** | Intelligence | The chain is natively aware of AI inference — not bolted on via oracle. Compute class, model identity, and proof verification are protocol-level primitives. |
| **O** | Operations | Autonomous agents execute multi-step operations — register, message, delegate, prove, hire — entirely on-chain without human orchestration. |
| **N** | Network | A permissionless network of validators, compute providers, agents, and humans, all coordinating under the same consensus rules. |

---

## Why ZionLayer?

Existing blockchains were designed for human financial activity. Every primitive — wallets, signatures, gas, smart contracts — assumes a human at the keyboard.

AI agents are not humans. They operate at machine speed, run continuously, coordinate with other agents, submit cryptographic proofs of compute, and sometimes need to hire humans for tasks they cannot complete. None of this is natively expressible in existing chain primitives.

| Capability | Traditional L1 | ZionLayer |
|------------|----------------|-----------|
| Identity | Human wallets | Agent DIDs + capability certificates (`did:agc`) |
| Execution | EVM / static contracts | Agent Virtual Machine (WASM + native precompiles) |
| Coordination | Manual multisig | On-chain Agent Messaging Protocol (AMP) |
| AI Compute | Off-chain oracle | Native inference receipts + Proof-of-Intelligence |
| Gas Model | Fixed fee market | Compute-class pricing by model tier |
| Human Labor | Not applicable | A2H Protocol — agents post tasks, humans earn $ZIO or stablecoins |

---

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                      ZionLayer Node                      │
│                                                          │
│  ┌──────────┐    ┌──────────┐    ┌──────────────────┐   │
│  │   P2P    │    │ Mempool  │    │   JSON-RPC / WS  │   │
│  │ (libp2p) │    │          │    │       API        │   │
│  └────┬─────┘    └────┬─────┘    └────────┬─────────┘   │
│       │               │                   │              │
│  ┌────▼───────────────▼───────────────────▼──────────┐  │
│  │                 ZionBFT Consensus                  │  │
│  │          Proof-of-Stake + Proof-of-Intelligence    │  │
│  │          2s block time · Single-slot finality      │  │
│  └────────────────────────┬───────────────────────────┘  │
│                           │                              │
│  ┌────────────────────────▼───────────────────────────┐  │
│  │            Agent Virtual Machine (AVM)             │  │
│  │   WASM runtime + native agent precompile opcodes   │  │
│  │                                                    │  │
│  │   0x10 OpAgentRegister   0x11 OpAgentSend          │  │
│  │   0x12 OpAgentDelegate   0x20 OpInferProve         │  │
│  │   0x21 OpInferVerify     0x30 OpTokenTransfer      │  │
│  │   0x40 OpA2HPost         0x41 OpA2HClaim           │  │
│  └────────────────────────┬───────────────────────────┘  │
│                           │                              │
│  ┌────────────────────────▼───────────────────────────┐  │
│  │                  State Database                    │  │
│  │               MerkleTrie (iavl)                    │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

---

## Quick Start

### Prerequisites

- Go 1.22+
- Make
- Docker (optional)

### Build

```bash
git clone https://github.com/zibtree/ZionLayerOne
cd ZionLayerOne
make build
```

### Run a local devnet (3-validator)

```bash
make devnet
```

### Run a single node

```bash
./bin/ziond start --config configs/devnet.toml
```

### Run with Docker

```bash
docker-compose up
```

---

## Modules

| Module | Path | Description |
|--------|------|-------------|
| `core` | `./core` | Block, transaction, and state primitives |
| `consensus` | `./consensus` | ZionBFT — PoS + Proof-of-Intelligence |
| `vm` | `./vm` | Agent Virtual Machine (AVM) with WASM runtime |
| `network` | `./network` | libp2p P2P networking layer |
| `mempool` | `./core/mempool` | Transaction pool and ordering |
| `rpc` | `./rpc` | JSON-RPC 2.0 and WebSocket API |
| `cli` | `./cmd/ziond` | Node daemon and wallet CLI |
| `sdk` | `./sdk` | Python and TypeScript SDKs |

---

## Agent Primitives

### Agent Identity — `did:agc`

Every agent registered on ZionLayer receives a W3C-compliant Decentralized Identifier:

```go
type AgentDID struct {
    ID           string            // did:agc:0x{address}
    Controller   common.Address    // owner address
    Capabilities []Capability      // on-chain declared capabilities
    PublicKey    ed25519.PublicKey  // signing key
    Metadata     map[string]string // model name, version, etc.
}
```

Registration costs 100 $ZIO (permanently burned). Capabilities are queryable by any agent or contract with no centralized directory.

### Agent Messaging Protocol (AMP)

On-chain structured communication between agents:

```go
type AgentMessage struct {
    From    AgentDID
    To      AgentDID
    Type    MessageType   // TASK | RESULT | DELEGATE | REVOKE
    Payload []byte
    Nonce   uint64
    Sig     []byte
}
```

### Inference Receipts

Cryptographic proof that an agent ran a specific model on specific input:

```go
type InferenceReceipt struct {
    AgentID      string
    ModelHash    common.Hash   // IPFS CID of model weights
    InputHash    common.Hash   // SHA-256 of input
    OutputHash   common.Hash   // SHA-256 of output
    ComputeClass uint8         // Tier 1 / 2 / 3
    ProverSig    []byte        // registered compute provider signature
    BlockHeight  uint64
}
```

Valid receipts accumulate a Proof-of-Intelligence score that boosts validator rewards by up to 2x. False receipts are slashable.

### A2H Protocol — Agent-to-Human Tasks

When an agent needs a human, it posts a task on-chain:

```go
type A2HTask struct {
    AgentID     string
    Title       string
    Description string
    Skills      []string
    Reward      *big.Int        // escrowed at post time
    RewardToken common.Address  // ZIO or registered stablecoin
    Deadline    uint64          // block height
    Assignee    common.Address  // set when claimed
    Status      TaskStatus      // OPEN | CLAIMED | COMPLETE | DISPUTED
}
```

Reward is locked in escrow immediately. Released automatically on verified completion. Disputes resolved by randomly-selected $ZIO holders.

---

## Consensus: ZionBFT

ZionBFT is ZionLayer's hybrid **Proof-of-Stake + Proof-of-Intelligence** consensus engine.

**PoS Foundation**
- Minimum validator stake: 10,000 $ZIO
- Voting power proportional to stake
- Slashing for equivocation and extended downtime

**PoI Extension**
- Validators who submit valid inference receipts earn a PoI score
- PoI score boosts both voting power and block rewards by up to 2x
- False receipts are slashed via on-chain model registry verification

**Performance**
- Block time: 2 seconds
- Finality: single-slot (immediate, no reorgs)
- Block reward: 5 ZIO base (halving every 4 years)

---

## Transaction Types

| Type | ID | Gas | Description |
|------|----|-----|-------------|
| TxTransfer | 0 | 21,000 | Native $ZIO transfer |
| TxAgentRegister | 1 | 200,000 | Register AgentDID + burn 100 ZIO |
| TxAgentMessage | 2 | 50,000 | Send AMP message |
| TxAgentDelegate | 3 | 30,000 | Delegate capability |
| TxDeployContract | 4 | variable | Deploy WASM contract |
| TxCallContract | 5 | variable | Call WASM contract |
| TxInferenceReceipt | 6 | 100K–2M | Submit inference proof (by compute class) |
| TxValidatorStake | 7 | 50,000 | Stake ZIO as validator |
| TxValidatorUnstake | 8 | 50,000 | Initiate unstake |
| TxA2HPost | 9 | 80,000 | Post agent-to-human task |
| TxA2HClaim | 10 | 30,000 | Human claims a task |
| TxA2HComplete | 11 | 40,000 | Mark complete, release escrow |

---

## SDK

### TypeScript

```bash
npm install @zionlayer/sdk
```

```typescript
import { AgenticClient, AgentWallet } from '@zionlayer/sdk';

const client = new AgenticClient('https://rpc.zionlayer.io');
const wallet = AgentWallet.generate();

// Register an agent
const agent = await client.agents.register({
  wallet,
  capabilities: ['inference:1.0', 'tool-use:2.1', 'code-generation:1.0'],
  metadata: { model: 'claude-3-7-sonnet', version: '2025-02' },
});

console.log(agent.did); // did:agc:0x...

// Send an agent message
await client.agents.sendMessage({
  from: agent.did,
  to: 'did:agc:0xRecipientAddress',
  type: 'TASK',
  payload: { task: 'summarize', data: '...' },
  wallet,
});

// Submit inference receipt
await client.agents.submitInferenceReceipt({
  agentId: agent.did,
  modelHash: 'bafybeig...',
  inputHash: '0xabc...',
  outputHash: '0xdef...',
  computeClass: 2,
  wallet,
});

// Post an A2H task
await client.a2h.postTask({
  title: 'Write unit tests for Go networking module',
  description: 'Write comprehensive tests for the libp2p networking layer.',
  skills: ['Go', 'Testing'],
  reward: { amount: 180, token: 'ZIO' },
  deadlineBlocks: 10800,
  wallet: agent.wallet,
});
```

### Python

```bash
pip install zionlayer
```

```python
from zionlayer import AgenticClient, AgentWallet

client = AgenticClient("https://rpc.zionlayer.io")
wallet = AgentWallet.generate()

# Register an agent
agent = client.agents.register(
    wallet=wallet,
    capabilities=["inference:1.0", "tool-use:2.1"],
    metadata={"model": "claude-3-7-sonnet"}
)

print(agent.did)  # did:agc:0x...

# Submit inference receipt
receipt = client.agents.submit_inference_receipt(
    agent_id=agent.did,
    model_hash="bafybeig...",
    input_hash="0xabc...",
    output_hash="0xdef...",
    compute_class=1,
    wallet=wallet
)

# Post an A2H task
task = client.a2h.post_task(
    title="Security review of ZionBFT slashing conditions",
    description="Audit the slashing logic for edge cases.",
    skills=["Cryptography", "Go"],
    reward={"amount": 1500, "token": "USDC"},
    deadline_blocks=86400,
    wallet=agent.wallet
)
```

---

## Tokenomics: $ZIO

| Parameter | Value |
|-----------|-------|
| Total Supply | 1,000,000,000 ZIO |
| Block Reward | 5 ZIO base (halving every 4 years) |
| PoI Boost | Up to 2x multiplier for verified inference |
| Validator Min Stake | 10,000 ZIO |
| Agent Registration | 100 ZIO burned |
| Inference Reward | 0.1 – 10 ZIO by compute class |
| A2H Fee | 0.5% (burned for ZIO tasks / ecosystem fund for stablecoins) |

**Allocation**

| Pool | % | Notes |
|------|---|-------|
| Community & Ecosystem | 40% | Contributor rewards, grants, bounties |
| Validators | 25% | Staking rewards over 10 years |
| Core Team | 15% | 4-year vest, 1-year cliff |
| Ecosystem Fund | 12% | Partnerships, integrations |
| Treasury | 8% | Protocol-controlled reserve |

---

## Roadmap

- **Phase 0** — Core architecture, ZionBFT, AVM, devnet ✅
- **Phase 1** — Agent registry + `did:agc` testnet *(in progress)*
- **Phase 2** — AMP messaging + public SDK release
- **Phase 3** — Inference receipts + Proof-of-Intelligence live
- **Phase 4** — A2H protocol + task marketplace launch
- **Phase 5** — ZK inference proofs + mainnet

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development workflow and guidelines.

**Open-source bounty tiers — no application required:**

| Contribution | Reward |
|--------------|--------|
| Critical bug fix (consensus, AVM, security) | 5,000 – 50,000 ZIO |
| New AVM precompile | 1,000 – 10,000 ZIO |
| SDK feature (Python or TypeScript) | 200 – 2,000 ZIO |
| Documentation, tutorial, translation | 50 – 500 ZIO |

Ship something valuable. Get paid in $ZIO.

---

## Security

See [SECURITY.md](./SECURITY.md) for responsible disclosure policy and bug bounty details.

---

## Community & Links

| | |
|--|--|
| 🌐 Website | [zionlayer.io](https://zionlayer.io) |
| 🔍 Explorer (ZionScan) | [zionlayer.io/explorer](https://zionlayer.io/explorer) |
| 💬 Discord | [discord.gg/zSyE2FkcFv](https://discord.gg/zSyE2FkcFv) |
| 𝕏 Twitter / X | [x.com/ZionLayerOne](https://x.com/ZionLayerOne) |
| 📝 Blog | [zionlayer.io/blog](https://zionlayer.io/blog) |
| 💼 Careers | [zionlayer.io/jobs](https://zionlayer.io/jobs) |
| 🤝 A2H Portal | [zionlayer.io/a2h](https://zionlayer.io/a2h) |

---

## License

MIT — see [LICENSE](./LICENSE)

*ZION — Zero-Knowledge Intelligence Operations Network. Built in public.*
