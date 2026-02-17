"""Pytest configuration and fixtures."""
import pytest
from pkg.collector.observer import RuntimeObserver
from pkg.emitter.publisher import EmissionPublisher
from pkg.models.emission import Emission, EmissionType

@pytest.fixture
def observer():
    """Create a test observer."""
    return RuntimeObserver()

@pytest.fixture
def publisher():
    """Create a test publisher."""
    return EmissionPublisher()

@pytest.fixture
def sample_emission():
    """Create a sample emission."""
    return Emission(
        id="em-test",
        emission_type=EmissionType.MODEL_USED,
        agent_id="agent-1",
        component_name="GPT-4",
        component_version="1.0",
        provider="OpenAI",
    )
