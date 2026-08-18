# swt3-ai-go

SWT3 AI Witness SDK for Go. Cryptographic attestation primitives for AI inference witnessing. Tamper-proof evidence that your AI is doing what you say it does.

**114 procedures | 62 namespaces | 10 SDKs | Zero dependencies**

Go is the dominant language for backend infrastructure at scale. Kubernetes operators, API gateways, data pipelines, and inference orchestrators are overwhelmingly written in Go. This package gives infrastructure teams the same cryptographic witnessing primitives available in the Python, TypeScript, Rust, C#, Ruby, Swift, and Kotlin SDKs, with zero external dependencies and byte-identical output across all 10 implementations.

EU AI Act GPAI transparency obligations enforced since **August 2, 2026**. High-risk enforcement follows **December 2, 2027**. This SDK gives you the cryptographic primitives for both.

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

## API

| Function | Description |
|----------|-------------|
| `MintFingerprint(tenant, proc, fa, fb, fc, ts)` | 12-char SHA-256 witness fingerprint |
| `MintLifecycleChainID(tenant, proc, fp, ts)` | `LC-` prefixed 16-char lifecycle chain ID |
| `SignPayload(key, fingerprint, agentID)` | HMAC-SHA256 payload signature |
| `Sha256Truncated(data, length)` | Truncated SHA-256 hex digest |
| `TimestampMs()` | Current time as (ms, epoch) |

## Cross-Language Parity

Output is byte-identical to the Python, TypeScript, Rust, C#, Ruby, Swift, and Kotlin SDKs. All 55 fingerprint vectors, 2 signing vectors, 5 hash vectors, and 3 lifecycle chain vectors pass.

## License

Apache-2.0
