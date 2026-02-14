# runtime-aibom-emitter

**Runtime AIBOM emitter** — generate an AI Bill of Materials from what *actually ran*, not just what was declared.

Combines runtime OTel trace evidence (models used, tools called, endpoints hit, frameworks loaded) with static inventory into a CycloneDX-compatible AIBOM. Ground truth for AI governance.

> Part of the **GenAI Infrastructure Standard** — a composable suite of open-source tools for enterprise-grade GenAI observability, security, and governance.
>
> | Layer | Component | Repo |
> |-------|-----------|------|
> | Privacy | Prompt Vault Processor | [prompt-vault-processor](https://github.com/nostalgicskinco/prompt-vault-processor) |
> | Normalization | Semantic Normalizer | [genai-semantic-normalizer](https://github.com/nostalgicskinco/genai-semantic-normalizer) |
> | Metrics | Cost & SLO Pack | [genai-cost-slo](https://github.com/nostalgicskinco/genai-cost-slo) |
> | Replay | Agent VCR | [agent-vcr](https://github.com/nostalgicskinco/agent-vcr) |
> | Testing | Regression Harness | [trace-regression-harness](https://github.com/nostalgicskinco/trace-regression-harness) |
> | Security | MCP Scanner | [mcp-security-scanner](https://github.com/nostalgicskinco/mcp-security-scanner) |
> | Gateway | MCP Policy Gateway | [mcp-policy-gateway](https://github.com/nostalgicskinco/mcp-policy-gateway) |
> | **Inventory** | **Runtime AIBOM Emitter** | **this repo** |

## Problem

Static AIBOMs (code scan) drift from reality. Your declared model list says "gpt-3.5-turbo" but runtime traces show the agent actually called gpt-4, claude-3-sonnet, and three undocumented tools. Compliance needs ground truth.

## Quick Start

```bash
# Build the CLI
go build -o aibom ./cmd/aibom

# Generate AIBOM from trace export
./aibom -trace traces.json -app "my-agent" -output aibom.json
```

## What Gets Extracted

| Component Type | Source Attributes |
|---------------|-------------------|
| Models | `gen_ai.request.model`, `llm.model_name`, `ai.model.id` |
| Providers | `gen_ai.system`, `llm.provider`, `ai.provider` |
| Tools | `tool.name`, `mcp.tool.name`, `gen_ai.tool.name` |
| Endpoints | `server.address`, `url.full`, `http.url` |
| Frameworks | OTel instrumentation scope name/version |

## Output Format

CycloneDX-AI compatible JSON with runtime evidence:

```json
{
  "bomFormat": "CycloneDX-AI",
  "specVersion": "1.6",
  "components": [...],
  "services": [...],
  "evidence": {
    "source": "runtime_traces",
    "traceIds": ["abc123"],
    "spanCount": 50,
    "findings": [...]
  }
}
```

## License

AGPL-3.0-or-later — see [LICENSE](LICENSE). Commercial licenses available — see [COMMERCIAL_LICENSE.md](COMMERCIAL_LICENSE.md).
