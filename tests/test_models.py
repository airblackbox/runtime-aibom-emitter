"""Test emission models."""
import pytest
from pkg.models.emission import EmissionType, Emission, EmissionSummary

def test_emission_type_enum():
    """Test EmissionType enum."""
    assert EmissionType.MODEL_USED.value == "model_used"
    assert EmissionType.TOOL_INVOKED.value == "tool_invoked"
    assert EmissionType.DATA_ACCESSED.value == "data_accessed"

def test_emission_creation():
    """Test Emission creation."""
    emission = Emission(
        emission_type=EmissionType.MODEL_USED,
        agent_id="agent-1",
        component_name="GPT-4",
        provider="OpenAI"
    )
    assert emission.agent_id == "agent-1"
    assert emission.component_name == "GPT-4"
    assert emission.provider == "OpenAI"
    assert emission.id == ""

def test_emission_with_id():
    """Test Emission with ID."""
    emission = Emission(
        id="em-123",
        emission_type=EmissionType.TOOL_INVOKED,
        agent_id="agent-1",
        component_name="SearchTool"
    )
    assert emission.id == "em-123"

def test_emission_summary():
    """Test EmissionSummary creation."""
    summary = EmissionSummary(
        agent_id="agent-1",
        total_emissions=10,
        unique_models=["GPT-4"],
        unique_tools=["SearchTool"],
    )
    assert summary.agent_id == "agent-1"
    assert summary.total_emissions == 10
    assert len(summary.unique_models) == 1
