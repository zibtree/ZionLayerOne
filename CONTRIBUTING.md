# Contributing to ZionLayer

**ZION — Zero-Knowledge Intelligence Operations Network**

**Protocol development status: active — September 2026**

ZionLayer is MIT licensed and built in public. Contributions are welcome across consensus, networking, execution, agent primitives, A2H, inference verification, SDKs, testing, and security.

## Current engineering priorities

1. Networked validator consensus
2. Deterministic execution and persistence
3. Agent capability delegation
4. A2H verification and dispute handling
5. Verifiable inference
6. Production-safe AVM/WASM execution
7. SDK correctness and integration tests

## Getting started

```bash
git clone https://github.com/zibtree/ZionLayerOne
cd ZionLayerOne
go version
make build
make test
```

## Development rules

- One logical change per PR.
- Include tests with protocol changes.
- Do not claim a feature is production-ready until implemented and tested.
- Update protocol documentation when public behavior changes.
- Never commit private keys, credentials, seed phrases, or production secrets.
- Security-sensitive changes require focused review before merge.

## Pull requests

Explain what changed, why it changed, security implications, compatibility implications, tests performed, and known limitations.

The default branch should remain releasable. Experimental protocol work belongs on feature branches until reviewed.

*ZionLayer — built in public for the agent economy.*