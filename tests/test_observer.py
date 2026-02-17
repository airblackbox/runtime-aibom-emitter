"""Test RuntimeObserver."""
import pytest
from pkg.collector.observer import RuntimeObserver
from pkg.models.emission import EmissionType

@pytest.mark.asyncio
async def test_observe_episodes():
    """Test observing episodes."""
    observer = RuntimeObserver()
    emissions = await observer.observe_episodes("agent-1", limit=5)
    assert len(emissions) > 0
    assert all(e.agent_id == "agent-1" for e in emissions)

@pytest.mark.asyncio
async def test_observe_multiple_times():
    """Test observing same agent multiple times."""
    observer = RuntimeObserver()
    first = await observer.observe_episodes("agent-1", limit=5)
    second = await observer.observe_episodes("agent-1", limit=5)
    assert len(second) == 0

@pytest.mark.asyncio
async def test_extract_emissions():
    """Test emission extraction."""
    observer = RuntimeObserver()
    episode_data = {
        "episode_id": "ep-1",
        "models_used": [
            {"name": "GPT-4", "version": "1.0", "provider": "OpenAI"}
        ],
        "tools_invoked": [],
        "data_accessed": [],
    }
    emissions = observer._extract_emissions("agent-1", episode_data)
    assert len(emissions) == 1
    assert emissions[0].emission_type == EmissionType.MODEL_USED
    assert emissions[0].component_name == "GPT-4"

def test_get_summary():
    """Test getting emission summary."""
    observer = RuntimeObserver()
    observer._emissions = [
        type('E', (), {
            'agent_id': 'agent-1',
            'emission_type': EmissionType.MODEL_USED,
            'component_name': 'GPT-4'
        })(),
        type('E', (), {
            'agent_id': 'agent-1',
            'emission_type': EmissionType.TOOL_INVOKED,
            'component_name': 'SearchTool'
        })(),
    ]
    summary = observer.get_summary("agent-1")
    assert summary.agent_id == "agent-1"
    assert summary.total_emissions == 2
