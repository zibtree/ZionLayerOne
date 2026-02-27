# Contributing to ZionLayer

Thanks for your interest in contributing! ZionLayer is open to contributions from everyone.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/zionlayer`
3. Install Go 1.22+
4. Run `make build` to verify your setup

## Development Workflow

```bash
# Run tests
make test

# Format code
make fmt

# Lint
make lint

# Start a local devnet
make devnet
```

## Areas to Contribute

### Core Protocol
- Improve ZionBFT consensus
- Optimize state trie performance
- Add P2P networking (libp2p integration)

### AVM Precompiles
- New agent primitive opcodes
- WASM contract support
- Inference receipt cryptographic verification

### SDKs
- Python SDK improvements
- TypeScript SDK improvements
- New language SDKs (Rust, Java, Go)

### Documentation
- Architecture deep-dives
- Tutorials and examples
- API reference

## Pull Request Guidelines

- Write clear commit messages
- Include tests for new functionality
- Update docs if you change public APIs
- Keep PRs focused — one change per PR

## Code Style

- Go: follow `gofmt` and `golangci-lint`
- TypeScript: follow `prettier` defaults
- Python: follow `black` formatting

## Questions?

Open an issue or join our Discord.
