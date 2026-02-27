/**
 * @zionlayer/sdk — TypeScript SDK for ZionLayer
 */

export interface AgentCapability {
  name: string;
  version: string;
}

export interface AgentDID {
  id: string;           // did:agc:0x...
  controller: string;   // owner address
  capabilities: AgentCapability[];
  publicKey: string;    // hex-encoded
  metadata?: Record<string, string>;
}

export type MessageType = 'TASK' | 'RESULT' | 'DELEGATE' | 'REVOKE';

export interface AgentMessage {
  from: string;
  to: string;
  type: MessageType;
  payload: unknown;
  nonce?: number;
}

export interface InferenceReceipt {
  agentId: string;
  modelHash: string;   // IPFS CID
  inputHash: string;
  outputHash: string;
  timestamp: number;
  proverSig: string;
}

export interface TransactionOptions {
  gasPrice?: bigint;
  gasLimit?: number;
  nonce?: number;
}

// ─── Client ────────────────────────────────────────────────────────────────

export class AgenticClient {
  readonly rpcUrl: string;
  readonly agents: AgentsAPI;
  readonly chain: ChainAPI;

  constructor(rpcUrl: string) {
    this.rpcUrl = rpcUrl;
    this.agents = new AgentsAPI(this);
    this.chain = new ChainAPI(this);
  }

  async call(method: string, params: unknown[]): Promise<unknown> {
    const res = await fetch(this.rpcUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        jsonrpc: '2.0',
        id: Date.now(),
        method,
        params,
      }),
    });
    const json = await res.json() as { result?: unknown; error?: { message: string } };
    if (json.error) throw new Error(json.error.message);
    return json.result;
  }
}

// ─── Agents API ────────────────────────────────────────────────────────────

export class AgentsAPI {
  constructor(private client: AgenticClient) {}

  /** Register a new AgentDID on-chain. */
  async register(
    wallet: AgentWallet,
    capabilities: AgentCapability[],
    metadata?: Record<string, string>
  ): Promise<{ did: AgentDID; txHash: string }> {
    const did: AgentDID = {
      id: `did:agc:${wallet.address}`,
      controller: wallet.address,
      capabilities,
      publicKey: wallet.publicKeyHex,
      metadata,
    };
    const tx = {
      type: 2, // TxAgentRegister
      from: wallet.address,
      gas: 200000,
      gasPrice: '1000000000',
      nonce: await this.getNonce(wallet.address),
      data: JSON.stringify(did),
    };
    const txHash = await this.client.call('zion_sendTransaction', [tx]) as string;
    return { did, txHash };
  }

  /** Fetch an agent record by DID. */
  async get(didId: string): Promise<AgentDID> {
    return this.client.call('zion_getAgent', [didId]) as Promise<AgentDID>;
  }

  /** Send an on-chain agent message. */
  async sendMessage(
    wallet: AgentWallet,
    msg: Omit<AgentMessage, 'nonce'>
  ): Promise<string> {
    const nonce = await this.getNonce(wallet.address);
    const fullMsg: AgentMessage = { ...msg, nonce };
    const tx = {
      type: 3, // TxAgentMessage
      from: wallet.address,
      gas: 50000,
      gasPrice: '1000000000',
      nonce,
      data: JSON.stringify(fullMsg),
    };
    return this.client.call('zion_sendTransaction', [tx]) as Promise<string>;
  }

  /** Submit an inference receipt for on-chain verification. */
  async submitInferenceReceipt(
    wallet: AgentWallet,
    receipt: InferenceReceipt
  ): Promise<string> {
    const tx = {
      type: 6, // TxInferenceReceipt
      from: wallet.address,
      gas: 100000,
      gasPrice: '1000000000',
      nonce: await this.getNonce(wallet.address),
      data: JSON.stringify(receipt),
    };
    return this.client.call('zion_sendTransaction', [tx]) as Promise<string>;
  }

  private async getNonce(address: string): Promise<number> {
    const acc = await this.client.call('zion_getBalance', [address]) as { nonce: string };
    return parseInt(acc.nonce ?? '0');
  }
}

// ─── Chain API ─────────────────────────────────────────────────────────────

export class ChainAPI {
  constructor(private client: AgenticClient) {}

  async getBalance(address: string): Promise<bigint> {
    const acc = await this.client.call('zion_getBalance', [address]) as { balance: string };
    return BigInt(acc.balance);
  }

  async getMempoolSize(): Promise<number> {
    const res = await this.client.call('zion_getMempoolSize', []) as { size: number };
    return res.size;
  }

  async getChainId(): Promise<string> {
    return this.client.call('zion_chainId', []) as Promise<string>;
  }
}

// ─── Wallet ────────────────────────────────────────────────────────────────

export class AgentWallet {
  readonly address: string;
  readonly publicKeyHex: string;
  private readonly privateKeyHex: string;

  private constructor(address: string, pubKey: string, privKey: string) {
    this.address = address;
    this.publicKeyHex = pubKey;
    this.privateKeyHex = privKey;
  }

  /** Generate a new random wallet (uses Web Crypto in browser, crypto in Node). */
  static generate(): AgentWallet {
    // In production: use noble-curves or ethers.js wallet
    const mockPriv = Array.from({ length: 32 }, () =>
      Math.floor(Math.random() * 256).toString(16).padStart(2, '0')
    ).join('');
    const mockPub = Array.from({ length: 33 }, () =>
      Math.floor(Math.random() * 256).toString(16).padStart(2, '0')
    ).join('');
    const addr = '0x' + mockPub.slice(-40);
    return new AgentWallet(addr, mockPub, mockPriv);
  }

  /** Load wallet from a hex private key. */
  static fromPrivateKey(privKeyHex: string): AgentWallet {
    // In production: derive public key and address from private key
    const addr = '0x' + privKeyHex.slice(-40);
    const pubKey = privKeyHex.padStart(66, '0');
    return new AgentWallet(addr, pubKey, privKeyHex);
  }
}

export default AgenticClient;
