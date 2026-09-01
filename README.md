Witness your AI. Prove it followed the rules. Cryptographic accountability for every inference, tool call, and resource access.

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://github.com/tenova-labs/swt3-ai-go/blob/main/LICENSE)

# swt3-ai-go

**SWT3 AI Witness SDK for Go**: mint, verify, and sign SWT3 witness anchors with cross-language parity. Zero external dependencies -- standard library only.

**266 procedures | 75 namespaces | 10 SDKs | Zero dependencies**

Go is the dominant language for backend infrastructure at scale. Kubernetes operators, API gateways, data pipelines, and inference orchestrators are overwhelmingly written in Go. This package gives infrastructure teams the same cryptographic witnessing primitives available in the Python, TypeScript, Rust, C#, Ruby, Swift, and Kotlin SDKs, with zero external dependencies and byte-identical output across all 10 implementations.

EU AI Act GPAI transparency obligations enforced since **August 2, 2026**. High-risk enforcement follows **December 2, 2027**. This SDK gives you the cryptographic primitives for both.

## What's New in v0.7.0

AI does not run in a vacuum. It authenticates with service accounts, runs on hardware with supply chains, and consumes enough electricity to reshape power grids. v0.7.0 extends the SWT3 protocol down the full AI infrastructure stack with 22 new procedures across three families: credential governance (NHI), hardware supply chains (HBOM/DPP), and energy management (ADR). Every new anchor type verifies with the existing `MintFingerprint` function -- no module update required.

**Why this matters for Go:** Go is where AI infrastructure lives. The Kubernetes operator scheduling GPU workloads, the API gateway routing inference requests, the data pipeline feeding training data -- all written in Go. These new procedure families map directly to what Go services manage. A Go-based orchestrator can now verify NHI credential rotation anchors, validate hardware inventory attestations, and confirm demand response settlement anchors -- all with the same `MintFingerprint` call. The fingerprint formula is unchanged. A credential rotation anchor and an inference anchor verify identically.

### 22 New Procedures

**NHI (Credential Governance) -- 6 procedures:** NHI-SCOPE.1, NHI-LIFE.1, NHI-PRIV.1, NHI-ROT.1, NHI-DEL.1, NHI-REV.1. Independent audit trail for every service account, API key, and machine credential your AI agents use. NIST 800-207 (Zero Trust), EU NIS-2 Art. 21, and CISA agentic AI guidance all require credential governance evidence.

**HBOM/DPP (Hardware + Battery Passport) -- 10 procedures:** HBOM-INV.1, HBOM-LIFE.1, HBOM-THERM.1, HBOM-WATER.1, HBOM-PUE.1, HBOM-SUPPLY.1, DPP-SOH.1, DPP-CHRG.1, DPP-DEGRAD.1, DPP-EOL.1. EU Battery Regulation (Feb 2027) digital passports, CRA hardware BOMs, and EED PUE reporting.

**ADR (Energy + Demand Response) -- 6 procedures:** ADR-EVENT.1, ADR-BASE.1, ADR-CURT.1, ADR-SETTLE.1, ADR-CARBON.1, ADR-GRID.1. Independent attestation for demand response events, baseline measurement, curtailment, settlement, carbon credits, and grid signal correlation.

### Updated Coverage

- 266 procedures across 75 namespaces (was 118/64)
- 59 MCP tools (was 37)
- 10 SDK languages with byte-identical output
- 36 framework crosswalks, 237 compliance guides
- ~2,950 tests passing across 5 languages

## What's New in v0.6.6

Supply chain accountability. Four new procedures, a CI/CD gate action, and OTel GenAI conventions across the ecosystem. Every improvement flows through to Go because fingerprints are identical across all 10 languages.

## What's New in v0.6.5

This is the debut release. Go is the 10th language in the SWT3 ecosystem, shipping alongside v0.6.5 "Scale Governance."

### Why Go, Why Now

**The infrastructure gap.** Python witnesses the model. TypeScript witnesses the API layer. Rust witnesses the inference engine. But the orchestration layer -- the Kubernetes operators that schedule workloads, the API gateways that route requests, the data pipelines that feed models -- is overwhelmingly Go. Without a Go SDK, infrastructure teams had to shell out to Python, use HTTP directly, or leave the orchestration layer unwitnessed. That gap is now closed.

**Why it matters:** When an assessor asks "was the AI monitored end-to-end?", a gap in the infrastructure layer means the answer is no. A Go API gateway that routes 10,000 inference requests per second can now mint fingerprints in-process with zero dependency overhead and sub-microsecond latency. The same fingerprint formula, the same signing algorithm, the same test vectors. No translation layer between what your gateway witnessed and what your Python model witnessed.

### What Ships

- `MintFingerprint` -- canonical 12-char SHA-256 fingerprint, byte-identical to all other SDKs
- `SignPayload` -- HMAC-SHA256 with agent identity binding
- `MintLifecycleChainID` -- lifecycle chain IDs for multi-anchor governance sequences
- `Sha256Truncated` -- content hashing for prompts, responses, model weights
- `TimestampMs` -- protocol-compatible millisecond timestamps
- All 65 test vectors passing (55 fingerprint + 2 signing + 5 hash + 3 lifecycle chain)
- Zero external dependencies (standard library only)

### What Else Shipped in v0.6.5

- **Probabilistic witnessing** (Python + TypeScript): sample inferences at configurable rates while keeping safety procedures at 100%. Deterministic hash-based sampling -- any Go verifier can reproduce the sampling decision.
- **Governance effectiveness metadata** (Python + TypeScript): review duration and participant count attached to governance anchors. Forward-compatible.
- **MCP Witness Middleware**: `withSWT3(transport)` wraps any MCP transport for zero-code tool call witnessing.
- **AI-SAMPLE.1**: new procedure for probabilistic witnessing summaries. 114 total procedures, 62 namespaces.
- **2,515 tests** passing across 5 languages. 36 framework crosswalks. 215 compliance guides.

## Install

```bash
go get github.com/tenova-labs/swt3-ai-go
```

## Quick Start

```go
package main

import (
    "fmt"
    swt3 "github.com/tenova-labs/swt3-ai-go"
)

func main() {
    // Mint a fingerprint
    ts, _ := swt3.TimestampMs()
    fp := swt3.MintFingerprint("MY_TENANT", "AI-INF.1", 1, 1, 0, ts)
    fmt.Println("Fingerprint:", fp)

    // Sign the anchor
    sig := swt3.SignPayload("my-signing-key", fp, "agent-1")
    fmt.Println("Signature:", sig)

    // Lifecycle chain ID
    chainID := swt3.MintLifecycleChainID("MY_TENANT", "AI-EMRG.1", fp, ts)
    fmt.Println("Chain ID:", chainID)

    // Content hash
    hash := swt3.Sha256Truncated("gpt-4o-2024-11-20", 16)
    fmt.Println("Model hash:", hash)
}
```

## What You Get

| Function | Description |
|----------|-------------|
| `MintFingerprint(tenant, proc, fa, fb, fc, ts)` | 12-char SHA-256 witness fingerprint |
| `MintLifecycleChainID(tenant, proc, fp, ts)` | `LC-` prefixed 16-char lifecycle chain ID |
| `SignPayload(key, fingerprint, agentID)` | HMAC-SHA256 payload signature |
| `Sha256Truncated(data, length)` | Truncated SHA-256 hex digest |
| `TimestampMs()` | Current time as (ms, epoch) |

All output is byte-identical to the Python, TypeScript, Swift, Rust, C#, Ruby, Kotlin, and MCP SDKs. 10 languages, one audit trail. Verified by shared test vectors at build time.

## Verify Any Anchor From Your Terminal

```bash
echo -n "WITNESS:DEMO_TENANT:AI-INF.1:1:1:0:1774800000000" | sha256sum | cut -c1-12
# Produces a 12-character fingerprint. Compare it to the anchor. If it matches, the anchor is real.
```

No SDK needed. Works on any machine, any language.

## Cross-Language Parity

All SWT3 SDKs produce identical fingerprints from the same inputs. A unified audit trail across your entire stack, verified by shared test vectors at build time.

| Language | Package | Registry |
|----------|---------|----------|
| Python | [swt3-ai](https://pypi.org/project/swt3-ai/) | PyPI |
| TypeScript | [@tenova/swt3-ai](https://www.npmjs.com/package/@tenova/swt3-ai) | npm |
| Swift | [swt3-ai](https://github.com/tenova-labs/swt3-ai-swift) | Swift Package Index |
| Rust | [swt3-ai](https://crates.io/crates/swt3-ai) | crates.io |
| C# / .NET | [swt3-ai](https://www.nuget.org/packages/swt3-ai) | NuGet |
| Ruby | [swt3-ai](https://rubygems.org/gems/swt3-ai) | RubyGems |
| Go | swt3-ai-go (this package) | Go modules |
| MCP Server | [@tenova/swt3-mcp](https://www.npmjs.com/package/@tenova/swt3-mcp) | npm + MCP Registry |
| K8s Witness Agent | [swt3-witness](https://github.com/tenova-labs/swt3-ai/tree/main/packages/swt3-witness) | GHCR + Helm |

The Python and TypeScript SDKs include the full witness pipeline: transparent client wrapping, buffer management, clearing engine, adapter support (OpenAI, Anthropic, Bedrock, vLLM, Ollama, LangChain, LangGraph, Microsoft AGT, Google ADK, CrewAI), trust mesh, policy-as-code, and Merkle accumulator. Use them for production AI witnessing. Use this Go package for embedding fingerprint minting and verification into infrastructure-layer code: API gateways, Kubernetes operators, data pipelines, or any Go service that touches AI.

## Regulatory Coverage

The SWT3 AI Witnessing Profile maps to:

- **EU AI Act**: Articles 9, 10, 12, 13, 14, 53, 72
- **NIST AI RMF**: GOVERN, MAP, MEASURE, MANAGE functions
- **ISO 42001**: Annex A AI management controls
- **NIST 800-53**: SI-7 (integrity), AU-2/AU-3 (audit), AC controls
- **SR 11-7**: Model risk management (financial services)

## Privacy

Your prompts and responses **never leave your infrastructure**. The SDK computes SHA-256 hashes locally and transmits only irreversible hashes and numeric factors. At Clearing Level 3, even the model name is hashed. The witness endpoint is a blind registrar: it stores cryptographic proofs, not your data.

## Links

- **Website**: [tenova.io](https://tenova.io)
- **Protocol Spec**: [SWT3-SPEC-v1.0](https://github.com/tenova-labs/swt3-ai)
- **Live Demo**: [sovereign.tenova.io/audit/axm_audit_demo_eu_ai_act_public](https://sovereign.tenova.io/audit/axm_audit_demo_eu_ai_act_public)

---

*SWT3: Sovereign Witness Traceability. We don't run your models. We witness them.*

SWT3 and Sovereign Witness Traceability are trademarks of Tenable Nova LLC. Patent pending. Apache 2.0 licensed.

This project is not affiliated with, endorsed by, or sponsored by any third-party AI provider. All third-party trademarks are the property of their respective owners. Use of these names is for identification and interoperability purposes only.
