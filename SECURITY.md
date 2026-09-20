# Security Policy

**ZionLayer security policy — September 2026**

ZionLayer is early-stage protocol infrastructure. The current repository is a development prototype and should **not** be used to custody production funds.

## Responsible disclosure

Please report vulnerabilities privately rather than opening a public issue.

**Email:** security@zionlayer.io

Include the vulnerability description, reproduction steps, affected component/version, impact assessment, proof of concept where safe, and suggested mitigation.

We aim to acknowledge reports within 48 hours and coordinate remediation based on severity.

## Priority security areas

- consensus safety and liveness
- double-spend prevention
- transaction authentication
- nonce/replay protection
- deterministic state transitions
- state persistence
- AVM sandbox escapes
- gas exhaustion
- A2H escrow integrity
- inference receipt forgery
- RPC denial of service
- SDK key handling

## Current limitations

Networked BFT quorum certificates, validator signature verification, slashing, production WASM sandboxing, cryptographic inference verification, and decentralized A2H dispute arbitration remain under development.

## Scope

In scope: ZionLayer node, consensus, state and transaction execution, AVM, A2H, JSON-RPC, and official SDKs.

Out of scope: third-party infrastructure not controlled by ZionLayer and purely theoretical reports without actionable impact.

*ZionLayer — security is a protocol requirement, not a marketing claim.*