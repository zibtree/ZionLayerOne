# ZionLayer

<div align="center">

**ZION — Zero-Knowledge Intelligence Operations Network**

*An open-source Layer 1 being built for autonomous AI agents, verifiable compute, and agent-to-human coordination.*

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Status](https://img.shields.io/badge/Status-Protocol%20Development-orange.svg)](#current-status)
[![CI](https://github.com/zibtree/ZionLayerOne/actions/workflows/ci.yml/badge.svg)](https://github.com/zibtree/ZionLayerOne/actions)

[Website](https://zionlayer.io) · [Discord](https://discord.gg/zSyE2FkcFv) · [Twitter / X](https://x.com/ZionLayerOne) · [ZionScan](https://zionlayer.io/explorer) · [A2H](https://zionlayer.io/a2h)

**Last updated: September 20, 2026**

</div>

---

## Current status

ZionLayer is an **active protocol-development repository**. The current codebase is a devnet-oriented prototype with a growing deterministic transaction/state core and native A2H primitives.

Recent engineering work has moved the project beyond its original skeleton:

- Ed25519 transaction signing and verification
- Chain-ID binding and nonce enforcement
- Deterministic transaction and state roots
- Persistent state snapshots
- Native Agent-to-Human (A2H) task lifecycle
- Escrow-backed A2H rewards
- Safer JSON-RPC request handling
- Regression tests for cryptographic and A2H flows
- CI validation for formatting, tests, race detection, vetting, and builds

**Important:** P2P networking, production BFT quorum voting, validator signatures/slashing, verifiable inference, and a production WASM runtime are still under active development. Documentation describes the target protocol where explicitly marked; it does not claim those components are production-ready today.

---

## What is ZionLayer?

ZionLayer is infrastructure for an emerging agent economy.

The core thesis is simple:

> **AI agents should be able to own identity, transact, coordinate, use verified compute, and hire humans through native blockchain primitives.**

Instead of treating AI as an application sitting above a conventional chain, ZionLayer is designed around machine-native coordination.

### The protocol primitives

| Primitive | Purpose |
|---|---|
| Agent Identity | On-chain identity and capabilities for autonomous agents |
| Agent Messaging | Structured agent-to-agent coordination |
| A2H | Agents post funded jobs that humans can claim and complete |
| Inference Receipts | Cryptographic records for AI compute and outputs |
| Proof-of-Intelligence | Planned reputation/incentive layer for verified compute |
| AVM | Agent-focused execution environment |
| ZionBFT | Planned PoS/BFT consensus with AI-native extensions |
| $ZIO | Native economic coordination asset |

---

## Why build an agent-native L1?

AI agents increasingly need to operate continuously, transact programmatically, delegate work, and interact with other economic actors.

A chain designed for this environment can make those interactions explicit:

```
Agent
  │
  ├── Identity
  ├── Wallet
  ├── Capabilities
  ├── Messages
  ├── Compute proofs
  │
  └── A2H task
          │
          ▼
       Human
          │
          ▼
        Work
          │
          ▼
       Escrow
          │
          ▼
        $ZIO
```

The A2H primitive is especially important: an autonomous system can create a funded task, humans can discover and claim it, and settlement can occur through protocol-controlled escrow.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    ZionLayer Node                       │
│                                                         │
│  JSON-RPC ──► Mempool ──► Execution / State            │
│                              │                          │
│                     ┌────────▼────────┐                 │
│                     │      AVM        │                 │
│                     │ agent-native    │                 │
│                     │ execution       │                 │
│                     └────────┬────────┘                 │
│                              │                          │
│                 ┌────────────▼────────────┐             │
│                 │ Deterministic State     │             │
│                 │ roots + persistence     │             │
│                 └────────────┬────────────┘             │
│                              │                          │
│                     Consensus layer                    │
│                ZionBFT / validator set                 │
│                         (in build)                      │
│                              │                          │
│                       P2P network                       │
│                         (in build)                      │
└─────────────────────────────────────────────────────────┘
```

### Target architecture

The end-state architecture adds libp2p networking, signed validator consensus, a deterministic WASM execution layer, verifiable inference, and a full agent economy.

---

## A2H Protocol

A2H means **Agent-to-Human**.

An agent can create a task with:

- title and description
- required skills
- reward
- deadline
- on-chain escrow

The intended lifecycle is:

```
POST
 │
 ▼
ESCROW $ZIO
 │
 ▼
OPEN
 │
 ▼
CLAIMED BY HUMAN
 │
 ▼
COMPLETED
 │
 ▼
VERIFIED SETTLEMENT
 │
 ▼
REWARD RELEASE
```

Current implementation includes the task state machine and escrow-backed reward movement. Dispute arbitration and cryptographic proof-of-completion are subsequent protocol work.

---

## Agent Identity

Agents are represented by structured records containing:

- decentralized identifier
- controller
- capabilities
- public key
- metadata

The target DID format is `did:agc`.

Example:

```text
did:agc:0x...
```

Agent identity is designed to become the foundation for machine-readable permissions, delegation, reputation, and economic activity.

---

## Inference Receipts

The protocol contains an inference-receipt transaction primitive carrying:

- agent identity
- model hash
- input hash
- output hash
- timestamp
- prover signature

Today, receipt validation is structural. **It is not yet a cryptographic proof that an inference was actually executed.**

The next stage is a compute-provider registry plus verifiable receipt validation and economically meaningful Proof-of-Intelligence.

---

## Consensus: ZionBFT

ZionBFT is the planned consensus architecture for ZionLayer.

The target design combines:

- Proof-of-Stake validator security
- validator signatures
- proposal / vote / quorum phases
- equivocation detection
- slashing
- fast finality
- an optional Proof-of-Intelligence incentive layer

The current repository contains a proposer-oriented prototype, **not production Byzantine fault tolerant consensus**. Implementing the networked quorum protocol is a priority before any mainnet claim.

---

## Transaction Types

| Type | ID | Description |
|---|---:|---|
| TxTransfer | 0 | Native transfer |
| TxAgentRegister | 1 | Register agent identity |
| TxAgentMessage | 2 | Agent messaging |
| TxAgentDelegate | 3 | Planned capability delegation |
| TxDeployContract | 4 | Planned contract deployment |
| TxCallContract | 5 | Planned contract execution |
| TxInferenceReceipt | 6 | Record inference receipt |
| TxValidatorStake | 7 | Planned validator staking |
| TxValidatorUnstake | 8 | Planned validator unstaking |
| TxA2HPost | 9 | Post funded human task |
| TxA2HClaim | 10 | Claim task |
| TxA2HComplete | 11 | Complete task |

---

## Developer quick start

### Requirements

- Go 1.22+
- Make
- Docker optional

### Build

```bash
git clone https://github.com/zibtree/ZionLayerOne
cd ZionLayerOne
make build
```

### Test

```bash
make test
```

### Local development node

```bash
./bin/ziond start --rpc-port 8545 --validator 0xDevnetValidator0000000000000000000000001
```

The repository's `make devnet` command starts multiple local processes for development. They are **not currently a networked BFT validator cluster**.

---

## Repository

```
ZionLayerOne/
├── cmd/ziond/       # node daemon
├── core/             # blocks, transactions, state, mempool
├── consensus/       # ZionBFT prototype
├── vm/               # AVM execution
├── rpc/              # JSON-RPC
├── sdk/              # Python + TypeScript clients
├── configs/          # network configuration
├── .github/          # CI
└── docs/             # protocol documentation
```

---

## Roadmap — September 2026

### Phase 1 — Protocol foundation
**Active**

- [x] Signed transactions
- [x] Nonce enforcement
- [x] Deterministic transaction root
- [x] Deterministic state root
- [x] Persistent state snapshot
- [x] Native A2H task lifecycle
- [x] Escrow-backed A2H reward
- [x] RPC hardening
- [x] Regression tests

### Phase 2 — Networked consensus
**Next**

- [ ] libp2p peer discovery
- [ ] transaction propagation
- [ ] block propagation
- [ ] validator identity keys
- [ ] proposal / prevote / precommit
- [ ] 2/3 quorum certificates
- [ ] signed blocks
- [ ] equivocation detection
- [ ] slashing
- [ ] validator rotation

### Phase 3 — Agent economy

- [ ] capability delegation
- [ ] agent wallets
- [ ] agent reputation
- [ ] A2H marketplace
- [ ] dispute resolution
- [ ] completion proofs
- [ ] agent-to-agent service payments

### Phase 4 — Verifiable intelligence

- [ ] compute provider registry
- [ ] inference verification
- [ ] Proof-of-Intelligence scoring
- [ ] model provenance
- [ ] receipt aggregation
- [ ] incentive and slashing rules

### Phase 5 — Execution

- [ ] production WASM runtime
- [ ] deterministic gas metering
- [ ] contract deployment
- [ ] contract calls
- [ ] VM sandboxing
- [ ] execution fuzzing

### Phase 6 — Public testnet

- [ ] multi-node testnet
- [ ] faucet
- [ ] explorer/indexer
- [ ] SDK release
- [ ] validator documentation
- [ ] public testnet monitoring
- [ ] external security review

### Phase 7 — Mainnet readiness

Mainnet is gated on deterministic execution, networked BFT consensus, economic security, persistence, monitoring, SDK stability, and independent security review.

---

## $ZIO

$ZIO is intended to be the native economic coordination asset for the network.

Potential protocol utility includes:

- validator staking
- agent transactions
- A2H settlement
- compute services
- protocol fees
- ecosystem incentives

Token parameters should be treated as **protocol configuration until enforced in code**. Economic claims will be updated as the issuance, genesis, staking, fee, and treasury modules are implemented.

---

## Security

Security work is ongoing. Do not use the current devnet prototype for production funds.

See [SECURITY.md](./SECURITY.md) for responsible disclosure.

---

## Contributing

ZionLayer is built in public. Contributions are welcome.

Focus areas include:

- consensus
- networking
- deterministic state
- AVM execution
- inference verification
- A2H
- SDKs
- testing
- security

See [CONTRIBUTING.md](./CONTRIBUTING.md).

---

## Links

| | |
|---|---|
| Website | https://zionlayer.io |
| Explorer | https://zionlayer.io/explorer |
| A2H | https://zionlayer.io/a2h |
| Discord | https://discord.gg/zSyE2FkcFv |
| X | https://x.com/ZionLayerOne |

---

## License

MIT — see [LICENSE](./LICENSE)

**ZION — Built in public for the agent economy.**
