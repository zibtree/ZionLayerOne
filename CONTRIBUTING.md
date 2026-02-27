# Contributing to ZionLayer

**ZION — Zero-Knowledge Intelligence Operations Network**

ZionLayer is MIT licensed and built entirely in public. Contributions from anyone are welcome — no application, no committee, no approval required. Ship something valuable and get paid in $ZIO from the ecosystem fund.

---

## Bounty Tiers

| Contribution | Reward |
|---|---|
| Critical bug fix (consensus, AVM, security) | 5,000 – 50,000 ZIO |
| New AVM precompile | 1,000 – 10,000 ZIO |
| SDK feature (Python or TypeScript) | 200 – 2,000 ZIO |
| Documentation, tutorial, translation | 50 – 500 ZIO |

Open an issue first for anything above 1,000 ZIO to align on scope before building.

---

## Getting Started

```bash
git clone https://github.com/zibtree/ZionLayerOne
cd ZionLayerOne
go version  # requires 1.22+
make build
make test
```

### Run a local devnet

```bash
make devnet
# Starts 3 validators on localhost, RPC at :8545
```

### Development commands

```bash
make build      # compile ziond binary
make test       # run full test suite
make fmt        # gofmt + goimports
make lint       # golangci-lint
make devnet     # 3-validator local network
make clean      # remove build artifacts
```

---

## Repository Structure

```
ZionLayerOne/
├── cmd/ziond/          # Node daemon entrypoint
├── core/               # Block, tx, state primitives
│   ├── block.go
│   ├── tx.go           # 12 transaction types
│   ├── statedb.go
│   └── mempool.go
├── consensus/          # ZionBFT engine
│   └── zionbft.go
├── vm/                 # Agent Virtual Machine
│   └── avm.go          # WASM runtime + agent precompiles
├── network/            # libp2p P2P layer
├── rpc/                # JSON-RPC 2.0 server
├── sdk/                # Python + TypeScript SDKs
│   ├── typescript/
│   └── python/
├── configs/            # devnet / testnet configs
├── docs/               # Protocol documentation
└── tests/              # Integration tests
```

---

## High-Priority Areas

### Core Protocol
- ZionBFT consensus optimizations and edge case handling
- State trie performance improvements
- P2P networking hardening
- Mempool ordering improvements

### AVM Precompiles
- New agent primitive opcodes (proposals welcome)
- WASM contract tooling and examples
- Inference receipt cryptographic verification hardening
- A2H task escrow and dispute resolution contracts

### SDKs
- Python SDK: async support, type hints, testing utilities
- TypeScript SDK: React hooks, browser wallet integration
- New language SDKs: Rust, Java, Go client library

### Documentation
- Architecture deep-dives
- Agent development tutorials
- A2H integration guide
- Validator setup guide

---

## Pull Request Guidelines

- Open an issue first for significant changes
- One logical change per PR — keep them focused
- Include tests for all new functionality
- Update documentation if you change public APIs
- Write clear commit messages: `type(scope): description`
  - Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

---

## Code Style

| Language | Formatter | Linter |
|---|---|---|
| Go | `gofmt`, `goimports` | `golangci-lint` |
| TypeScript | `prettier` | `eslint` |
| Python | `black` | `ruff` |

CI enforces all of the above. Run `make fmt && make lint` before opening a PR.

---

## Testing

```bash
# Unit tests
make test

# Integration tests (requires devnet running)
make devnet &
make test-integration

# Specific package
go test ./consensus/... -v
go test ./vm/... -v
```

---

## Questions

- Open a GitHub Discussion for architecture questions
- Join [Discord](https://discord.gg/zionlayer) for real-time conversation
- Email [dev@zionlayer.io](mailto:dev@zionlayer.io) for sensitive matters

---

*ZION — Zero-Knowledge Intelligence Operations Network. MIT License.*
