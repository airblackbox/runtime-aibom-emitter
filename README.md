# Runtime AIBOM Emitter

Generates AIBOM entries in real-time by observing live agent runs. Connects to Episode Store and Gateway to discover what models, tools, and data sources an agent actually uses during execution.

## Features

- Real-time emission tracking during agent execution
- Episode stream observation and parsing
- Integration with AIBOM Policy Engine
- Model/tool/data source usage discovery
- Emission filtering and summarization
- JSON export for audit trails

## Quick Start

```bash
pip install -e .
python -m app.server
```

API runs on `http://localhost:8700/v1`

## Emission Types

- **MODEL_USED**: Agent invoked an LLM model
- **TOOL_INVOKED**: Agent called a tool/function
- **DATA_ACCESSED**: Agent queried a data source
- **POLICY_APPLIED**: A policy was enforced

## API Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/v1/health` | Health check |
| POST | `/v1/emit` | Record emission |
| POST | `/v1/observe` | Observe episodes |
| GET | `/v1/summary/{agent_id}` | Get agent summary |
| GET | `/v1/emissions` | List emissions |
| POST | `/v1/publish` | Publish to AIBOM |
| POST | `/v1/export` | Export as JSON |

## Integration

Connect to Episode Store for live observation:
```
Episode Store (port 8000)
        ↓
    Observer
        ↓
  Emissions
        ↓
AIBOM Engine (port 8600)
```

## Testing

```bash
pytest tests/ -v
```

## License

MIT
