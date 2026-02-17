"""Test EmissionPublisher."""
import pytest
from pkg.emitter.publisher import EmissionPublisher
from pkg.models.emission import Emission, EmissionType

def test_publisher_initialization():
    """Test publisher initialization."""
    pub = EmissionPublisher()
    assert pub.aibom_api_url == "http://localhost:8600/v1"
    assert len(pub._emissions) == 0

def test_collect_emission(sample_emission):
    """Test collecting an emission."""
    pub = EmissionPublisher()
    pub.collect(sample_emission)
    assert len(pub._emissions) == 1

def test_collect_batch():
    """Test collecting batch of emissions."""
    pub = EmissionPublisher()
    emissions = [
        Emission(
            emission_type=EmissionType.MODEL_USED,
            agent_id="agent-1",
            component_name="Model1"
        ),
        Emission(
            emission_type=EmissionType.TOOL_INVOKED,
            agent_id="agent-1",
            component_name="Tool1"
        ),
    ]
    pub.collect_batch(emissions)
    assert len(pub._emissions) == 2

@pytest.mark.asyncio
async def test_publish():
    """Test publishing to AIBOM."""
    pub = EmissionPublisher()
    emission = Emission(
        id="em-1",
        emission_type=EmissionType.MODEL_USED,
        agent_id="agent-1",
        component_name="GPT-4",
        provider="OpenAI"
    )
    pub.collect(emission)
    result = await pub.publish("aibom-1")
    assert result["published"] is True
    assert result["count"] == 1

def test_get_emissions(sample_emission):
    """Test getting emissions."""
    pub = EmissionPublisher()
    pub.collect(sample_emission)
    emissions = pub.get_emissions()
    assert len(emissions) == 1

def test_get_emissions_filter_by_agent():
    """Test filtering emissions by agent."""
    pub = EmissionPublisher()
    pub.collect(
        Emission(
            emission_type=EmissionType.MODEL_USED,
            agent_id="agent-1",
            component_name="m1"
        )
    )
    pub.collect(
        Emission(
            emission_type=EmissionType.MODEL_USED,
            agent_id="agent-2",
            component_name="m2"
        )
    )
    emissions = pub.get_emissions(agent_id="agent-1")
    assert len(emissions) == 1
