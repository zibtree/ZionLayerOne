# Security Policy

**ZION — Zero-Knowledge Intelligence Operations Network**

ZionLayer is infrastructure. Security is not a feature — it is the foundation. If you discover a vulnerability, please do not open a public GitHub issue.

---

## Responsible Disclosure

**Email:** security@zionlayer.io

Include in your report:
- Clear description of the vulnerability
- Steps to reproduce
- Affected components and versions
- Potential impact assessment
- Suggested fix (if you have one)

We will acknowledge within **48 hours** and aim to resolve critical issues within **7 days**.

We will not pursue legal action against researchers who disclose responsibly and in good faith.

---

## Bug Bounty

| Severity | Examples | Reward |
|---|---|---|
| **Critical** | Consensus break, double-spend, fund theft, AVM escape | 50,000 – 200,000 ZIO |
| **High** | Node crash, state corruption, validator slashing bypass | 10,000 – 50,000 ZIO |
| **Medium** | DoS vector, mempool manipulation, RPC abuse | 1,000 – 10,000 ZIO |
| **Low** | Minor issues, information disclosure | 100 – 1,000 ZIO |

Bounties are paid in $ZIO from the treasury. We reserve the right to adjust rewards based on impact and quality of report.

---

## Scope

**In scope:**
- ZionLayer node (`ziond`)
- ZionBFT consensus logic
- Agent Virtual Machine (AVM) and precompiles
- A2H escrow and arbitration contracts
- JSON-RPC and WebSocket API
- Python SDK (`zionlayer`)
- TypeScript SDK (`@zionlayer/sdk`)

**Out of scope:**
- ZionScan front-end UI cosmetic bugs
- Third-party dependencies (report upstream)
- Theoretical attacks without a working proof-of-concept
- Issues in infrastructure not operated by ZionLayer

---

## Disclosure Timeline

| Stage | Target |
|---|---|
| Acknowledgement | 48 hours |
| Initial assessment | 5 business days |
| Fix development | Critical: 7 days / High: 30 days / Medium: 90 days |
| Public disclosure | Coordinated with reporter |

---

## PGP Key

For sensitive reports, encrypt your email to:

```
security@zionlayer.io
```

PGP key available on request via Discord or email.

---

*ZION — Zero-Knowledge Intelligence Operations Network. MIT License.*
