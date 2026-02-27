"""
zionlayer — Python SDK for ZionLayer
"""
from __future__ import annotations

import json
import secrets
import time
from dataclasses import dataclass, field, asdict
from typing import Any, Optional
from urllib.request import Request, urlopen


# ─── Data models ─────────────────────────────────────────────────────────────

@dataclass
class AgentCapability:
    name: str
    version: str


@dataclass
class AgentDID:
    id: str                                   # did:agc:0x...
    controller: str                            # owner address
    capabilities: list[AgentCapability]
    public_key: str                            # hex
    metadata: dict[str, str] = field(default_factory=dict)


@dataclass
class AgentMessage:
    from_did: str
    to_did: str
    type: str                                  # TASK | RESULT | DELEGATE | REVOKE
    payload: Any
    nonce: int = 0


@dataclass
class InferenceReceipt:
    agent_id: str
    model_hash: str                            # IPFS CID
    input_hash: str
    output_hash: str
    timestamp: int = field(default_factory=lambda: int(time.time()))
    prover_sig: str = ""


# ─── Wallet ───────────────────────────────────────────────────────────────────

class AgentWallet:
    def __init__(self, address: str, public_key: str, private_key: str):
        self.address = address
        self.public_key = public_key
        self._private_key = private_key

    @classmethod
    def generate(cls) -> "AgentWallet":
        """Generate a new random wallet."""
        priv = secrets.token_hex(32)
        pub = secrets.token_hex(33)
        addr = "0x" + pub[-40:]
        return cls(addr, pub, priv)

    @classmethod
    def from_private_key(cls, private_key_hex: str) -> "AgentWallet":
        """Load wallet from private key hex."""
        addr = "0x" + private_key_hex[-40:]
        pub = private_key_hex.zfill(66)
        return cls(addr, pub, private_key_hex)


# ─── Client ───────────────────────────────────────────────────────────────────

class AgenticClient:
    def __init__(self, rpc_url: str):
        self.rpc_url = rpc_url
        self.agents = AgentsAPI(self)
        self.chain = ChainAPI(self)

    def call(self, method: str, params: list[Any]) -> Any:
        payload = json.dumps({
            "jsonrpc": "2.0",
            "id": int(time.time() * 1000),
            "method": method,
            "params": params,
        }).encode()
        req = Request(self.rpc_url, data=payload, headers={"Content-Type": "application/json"})
        with urlopen(req) as resp:
            data = json.loads(resp.read())
        if "error" in data:
            raise RuntimeError(data["error"]["message"])
        return data.get("result")


# ─── Agents API ───────────────────────────────────────────────────────────────

class AgentsAPI:
    def __init__(self, client: AgenticClient):
        self._client = client

    def register(
        self,
        wallet: AgentWallet,
        capabilities: list[AgentCapability],
        metadata: Optional[dict[str, str]] = None,
    ) -> tuple[AgentDID, str]:
        """Register a new AgentDID on-chain. Returns (did, tx_hash)."""
        did = AgentDID(
            id=f"did:agc:{wallet.address}",
            controller=wallet.address,
            capabilities=capabilities,
            public_key=wallet.public_key,
            metadata=metadata or {},
        )
        tx = {
            "type": 2,
            "from": wallet.address,
            "gas": 200000,
            "gasPrice": "1000000000",
            "nonce": self._nonce(wallet.address),
            "data": json.dumps(asdict(did)),
        }
        tx_hash = self._client.call("zion_sendTransaction", [tx])
        return did, tx_hash

    def get(self, did_id: str) -> dict:
        """Fetch an agent record by DID string."""
        return self._client.call("zion_getAgent", [did_id])

    def send_message(self, wallet: AgentWallet, msg: AgentMessage) -> str:
        """Send an on-chain agent message. Returns tx hash."""
        msg.nonce = self._nonce(wallet.address)
        payload = asdict(msg)
        tx = {
            "type": 3,
            "from": wallet.address,
            "gas": 50000,
            "gasPrice": "1000000000",
            "nonce": msg.nonce,
            "data": json.dumps(payload),
        }
        return self._client.call("zion_sendTransaction", [tx])

    def submit_inference_receipt(
        self, wallet: AgentWallet, receipt: InferenceReceipt
    ) -> str:
        """Submit an inference receipt for on-chain verification."""
        tx = {
            "type": 6,
            "from": wallet.address,
            "gas": 100000,
            "gasPrice": "1000000000",
            "nonce": self._nonce(wallet.address),
            "data": json.dumps(asdict(receipt)),
        }
        return self._client.call("zion_sendTransaction", [tx])

    def _nonce(self, address: str) -> int:
        acc = self._client.call("zion_getBalance", [address]) or {}
        return int(acc.get("nonce", 0))


# ─── Chain API ────────────────────────────────────────────────────────────────

class ChainAPI:
    def __init__(self, client: AgenticClient):
        self._client = client

    def get_balance(self, address: str) -> int:
        acc = self._client.call("zion_getBalance", [address]) or {}
        return int(acc.get("balance", 0))

    def get_mempool_size(self) -> int:
        res = self._client.call("zion_getMempoolSize", []) or {}
        return res.get("size", 0)

    def get_chain_id(self) -> str:
        return self._client.call("zion_chainId", [])


# ─── Quick usage example ──────────────────────────────────────────────────────

if __name__ == "__main__":
    client = AgenticClient("http://localhost:8545")
    wallet = AgentWallet.generate()
    print(f"Wallet address: {wallet.address}")

    did, tx_hash = client.agents.register(
        wallet=wallet,
        capabilities=[
            AgentCapability("inference", "1.0"),
            AgentCapability("tool-use", "1.0"),
        ],
        metadata={"model": "claude-3-5-sonnet"},
    )
    print(f"Agent registered: {did.id}")
    print(f"TX hash: {tx_hash}")
