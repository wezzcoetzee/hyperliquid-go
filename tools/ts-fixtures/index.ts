// Captures golden signing fixtures from @nktkas/hyperliquid for Go parity tests.
// Pinned upstream version: @nktkas/hyperliquid@0.32.2 (see package.json).
//
// Public helpers used (from @nktkas/hyperliquid/signing):
//   - createL1ActionHash(action, nonce, vaultAddress?, expiresAfter?) -> 0x-hex keccak hash
//   - signL1Action({ wallet, action, nonce, isTestnet?, vaultAddress?, expiresAfter? })
//   - signUserSignedAction({ wallet, action, types })
//
// L1 action hash recipe (matches Hyperliquid spec):
//   keccak256( msgpack(action) || nonce_be8 || (vault ? 0x01 || addr20 : 0x00) || (expires ? 0x00 || expires_be8 : ε) )

import { writeFileSync, mkdirSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import { privateKeyToAccount } from "viem/accounts";
import {
  createL1ActionHash,
  signL1Action,
  signUserSignedAction,
} from "@nktkas/hyperliquid/signing";

const __dirname = dirname(fileURLToPath(import.meta.url));
const OUT_DIR = join(__dirname, "..", "..", "exchange", "testdata", "fixtures");

const PK = "0x0000000000000000000000000000000000000000000000000000000000000001" as const;
const NONCE = 1700000000000;
const wallet = privateKeyToAccount(PK);

mkdirSync(OUT_DIR, { recursive: true });

interface Fixture {
  privateKey: string;
  nonce: string;
  action: Record<string, unknown>;
  actionHash: string | null;
  signature: { r: string; s: string; v: number };
  userSigned: { type: string; chainId: number } | null;
}

function writeFixture(name: string, fx: Fixture) {
  const path = join(OUT_DIR, `${name}.json`);
  writeFileSync(path, JSON.stringify(fx, null, 2) + "\n");
  console.log("wrote", path);
}

// ---- order_l1 -----------------------------------------------------------
{
  const action = {
    type: "order",
    orders: [
      { a: 0, b: true, p: "30000", s: "0.1", r: false, t: { limit: { tif: "Gtc" } } },
    ],
    grouping: "na",
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("order_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- cancel_l1 ----------------------------------------------------------
{
  const action = {
    type: "cancel",
    cancels: [{ a: 0, o: 12345 }],
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("cancel_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- usd_send (user-signed EIP-712) -------------------------------------
{
  const action = {
    type: "usdSend",
    signatureChainId: "0xa4b1",
    hyperliquidChain: "Mainnet",
    destination: "0x0000000000000000000000000000000000000002",
    amount: "10",
    time: NONCE,
  } as const;
  const types = {
    "HyperliquidTransaction:UsdSend": [
      { name: "hyperliquidChain", type: "string" },
      { name: "destination", type: "string" },
      { name: "amount", type: "string" },
      { name: "time", type: "uint64" },
    ],
  };
  const signature = await signUserSignedAction({ wallet, action, types });
  writeFixture("usd_send", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash: null,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: { type: "HyperliquidTransaction:UsdSend", chainId: 0xa4b1 },
  });
}

console.log("done.");
