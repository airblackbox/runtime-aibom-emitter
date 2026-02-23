"""Emission models for runtime observation."""
from __future__ import annotations
from datetime import datetime, timezone
from enum import Enum
from typing import Any
from pydantic import BaseModel, Field

class EmissionType(str, Enum):
    """Types of emissions during agent execution."""
    MODEL_USED = "model_used"
    TOOL_INVOKED = "tool_invoked"
    DATA_ACCESSED = "data_accessed"
    POLICY_APPLIED = "policy_applied"

class Emission(BaseModel):
    """Single emission event during agent execution."""
    id: str = ""
    emission_type: EmissionType
    agent_id: str = ""
    timestamp: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    component_name: str = ""
    component_version: str = ""
    provider: str = ""
    metadata: dict[str, Any] = Field(default_factory=dict)

class EmissionSummary(BaseModel):
    """Summary of emissions for an agent."""
    agent_id: str
    total_emissions: int = 0
    unique_models: list[str] = Field(default_factory=list)
    unique_tools: list[str] = Field(default_factory=list)
    unique_data_sources: list[str] = Field(default_factory=list)
    observation_window_hours: float = 24.0
