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

// ---- order_l1_testnet ---------------------------------------------------
{
  const action = {
    type: "order",
    orders: [
      { a: 0, b: true, p: "30000", s: "0.1", r: false, t: { limit: { tif: "Gtc" } } },
    ],
    grouping: "na",
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: true });
  writeFixture("order_l1_testnet", {
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

// ---- cancel_by_cloid_l1 -------------------------------------------------
{
  const action = {
    type: "cancelByCloid",
    cancels: [{ asset: 0, cloid: "0xdeadbeef00000000000000000000000000000000000000000000000000000001" }],
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("cancel_by_cloid_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- modify_l1 ----------------------------------------------------------
{
  const action = {
    type: "modify",
    oid: 999,
    order: { a: 0, b: true, p: "31000", s: "0.2", r: false, t: { limit: { tif: "Gtc" } } },
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("modify_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- batch_modify_l1 ----------------------------------------------------
{
  const action = {
    type: "batchModify",
    modifies: [
      { oid: 111, order: { a: 0, b: true, p: "31000", s: "0.2", r: false, t: { limit: { tif: "Gtc" } } } },
      { oid: 222, order: { a: 1, b: false, p: "2000", s: "1.5", r: true, t: { limit: { tif: "Ioc" } } } },
    ],
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("batch_modify_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- schedule_cancel_l1 -------------------------------------------------
{
  const action = {
    type: "scheduleCancel",
    time: 1700000060000,
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("schedule_cancel_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- update_leverage_l1 -------------------------------------------------
{
  const action = {
    type: "updateLeverage",
    asset: 0,
    isCross: true,
    leverage: 10,
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("update_leverage_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- update_isolated_margin_l1 ------------------------------------------
{
  const action = {
    type: "updateIsolatedMargin",
    asset: 0,
    isBuy: true,
    ntli: 1000000,
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("update_isolated_margin_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- twap_order_l1 ------------------------------------------------------
{
  const action = {
    type: "twapOrder",
    twap: { a: 0, b: true, s: "1.0", r: false, m: 5, t: false },
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("twap_order_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- twap_cancel_l1 -----------------------------------------------------
{
  const action = {
    type: "twapCancel",
    a: 0,
    t: 42,
  } as const;
  const actionHash = createL1ActionHash({ action, nonce: NONCE });
  const signature = await signL1Action({ wallet, action, nonce: NONCE, isTestnet: false });
  writeFixture("twap_cancel_l1", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: null,
  });
}

// ---- withdraw3 (user-signed) --------------------------------------------
{
  const action = {
    type: "withdraw3",
    signatureChainId: "0xa4b1",
    hyperliquidChain: "Mainnet",
    destination: "0x0000000000000000000000000000000000000002",
    amount: "5",
    time: NONCE,
  } as const;
  const types = {
    "HyperliquidTransaction:Withdraw": [
      { name: "hyperliquidChain", type: "string" },
      { name: "destination", type: "string" },
      { name: "amount", type: "string" },
      { name: "time", type: "uint64" },
    ],
  };
  const signature = await signUserSignedAction({ wallet, action, types });
  writeFixture("withdraw3", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash: null,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: { type: "HyperliquidTransaction:Withdraw", chainId: 0xa4b1 },
  });
}

// ---- spot_send (user-signed) --------------------------------------------
{
  const action = {
    type: "spotSend",
    signatureChainId: "0xa4b1",
    hyperliquidChain: "Mainnet",
    destination: "0x0000000000000000000000000000000000000002",
    token: "USDC:0xeb62eee3685fc4c43992febcd9e75443",
    amount: "1",
    time: NONCE,
  } as const;
  const types = {
    "HyperliquidTransaction:SpotSend": [
      { name: "hyperliquidChain", type: "string" },
      { name: "destination", type: "string" },
      { name: "token", type: "string" },
      { name: "amount", type: "string" },
      { name: "time", type: "uint64" },
    ],
  };
  const signature = await signUserSignedAction({ wallet, action, types });
  writeFixture("spot_send", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash: null,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: { type: "HyperliquidTransaction:SpotSend", chainId: 0xa4b1 },
  });
}

// ---- usd_class_transfer (user-signed) -----------------------------------
{
  const action = {
    type: "usdClassTransfer",
    signatureChainId: "0xa4b1",
    hyperliquidChain: "Mainnet",
    amount: "100",
    toPerp: true,
    nonce: NONCE,
  } as const;
  const types = {
    "HyperliquidTransaction:UsdClassTransfer": [
      { name: "hyperliquidChain", type: "string" },
      { name: "amount", type: "string" },
      { name: "toPerp", type: "bool" },
      { name: "nonce", type: "uint64" },
    ],
  };
  const signature = await signUserSignedAction({ wallet, action, types });
  writeFixture("usd_class_transfer", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash: null,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: { type: "HyperliquidTransaction:UsdClassTransfer", chainId: 0xa4b1 },
  });
}

// ---- approve_agent (user-signed) ----------------------------------------
{
  const action = {
    type: "approveAgent",
    signatureChainId: "0xa4b1",
    hyperliquidChain: "Mainnet",
    agentAddress: "0x0000000000000000000000000000000000000003",
    agentName: "TestAgent",
    nonce: NONCE,
  } as const;
  const types = {
    "HyperliquidTransaction:ApproveAgent": [
      { name: "hyperliquidChain", type: "string" },
      { name: "agentAddress", type: "address" },
      { name: "agentName", type: "string" },
      { name: "nonce", type: "uint64" },
    ],
  };
  const signature = await signUserSignedAction({ wallet, action, types });
  writeFixture("approve_agent", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash: null,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: { type: "HyperliquidTransaction:ApproveAgent", chainId: 0xa4b1 },
  });
}

// ---- approve_builder_fee (user-signed) ----------------------------------
{
  const action = {
    type: "approveBuilderFee",
    signatureChainId: "0xa4b1",
    hyperliquidChain: "Mainnet",
    maxFeeRate: "0.1%",
    builder: "0x0000000000000000000000000000000000000004",
    nonce: NONCE,
  } as const;
  const types = {
    "HyperliquidTransaction:ApproveBuilderFee": [
      { name: "hyperliquidChain", type: "string" },
      { name: "maxFeeRate", type: "string" },
      { name: "builder", type: "address" },
      { name: "nonce", type: "uint64" },
    ],
  };
  const signature = await signUserSignedAction({ wallet, action, types });
  writeFixture("approve_builder_fee", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash: null,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: { type: "HyperliquidTransaction:ApproveBuilderFee", chainId: 0xa4b1 },
  });
}

// ---- token_delegate (user-signed) ---------------------------------------
{
  const action = {
    type: "tokenDelegate",
    signatureChainId: "0xa4b1",
    hyperliquidChain: "Mainnet",
    validator: "0x0000000000000000000000000000000000000005",
    wei: 1000000000000000000,
    isUndelegate: false,
    nonce: NONCE,
  } as const;
  const types = {
    "HyperliquidTransaction:TokenDelegate": [
      { name: "hyperliquidChain", type: "string" },
      { name: "validator", type: "address" },
      { name: "wei", type: "uint64" },
      { name: "isUndelegate", type: "bool" },
      { name: "nonce", type: "uint64" },
    ],
  };
  const signature = await signUserSignedAction({ wallet, action, types });
  writeFixture("token_delegate", {
    privateKey: PK,
    nonce: String(NONCE),
    action: action as unknown as Record<string, unknown>,
    actionHash: null,
    signature: { r: signature.r, s: signature.s, v: signature.v },
    userSigned: { type: "HyperliquidTransaction:TokenDelegate", chainId: 0xa4b1 },
  });
}

console.log("done.");
