"""Emission publisher for AIBOM integration."""
from __future__ import annotations
import json
from typing import Any
from pkg.models.emission import Emission, EmissionType

class EmissionPublisher:
    """Publishes emissions to AIBOM engine."""
    def __init__(self, aibom_api_url: str = "http://localhost:8600/v1") -> None:
        self.aibom_api_url = aibom_api_url
        self._emissions: list[Emission] = []
        self._published: set[str] = set()

    def collect(self, emission: Emission) -> None:
        """Collect an emission."""
        self._emissions.append(emission)

    def collect_batch(self, emissions: list[Emission]) -> None:
        """Collect multiple emissions."""
        self._emissions.extend(emissions)

    async def publish(
        self,
        aibom_id: str,
        agent_id: str | None = None
    ) -> dict[str, Any]:
        """Publish emissions to AIBOM engine."""
        target_emissions = self._emissions
        if agent_id:
            target_emissions = [
                e for e in self._emissions if e.agent_id == agent_id
            ]
        published_count = 0
        for emission in target_emissions:
            if emission.id not in self._published:
                component_type = self._map_emission_type(emission.emission_type)
                payload = {
                    "aibom_id": aibom_id,
                    "name": emission.component_name,
                    "component_type": component_type,
                    "provider": emission.provider,
                    "version": emission.component_version,
                    "description": f"Emitted from {emission.agent_id}",
                }
                self._published.add(emission.id)
                published_count += 1
        return {
            "published": True,
            "count": published_count,
            "aibom_id": aibom_id,
        }

    def export_json(self, filepath: str) -> None:
        """Export emissions to JSON file."""
        data = {
            "emissions": [e.dict() for e in self._emissions],
            "count": len(self._emissions),
        }
        with open(filepath, "w") as f:
            json.dump(data, f, indent=2, default=str)

    def get_emissions(
        self,
        agent_id: str | None = None,
        emission_type: EmissionType | None = None
    ) -> list[Emission]:
        """Get collected emissions with optional filtering."""
        result = self._emissions
        if agent_id:
            result = [e for e in result if e.agent_id == agent_id]
        if emission_type:
            result = [e for e in result if e.emission_type == emission_type]
        return result

    @staticmethod
    def _map_emission_type(emission_type: EmissionType) -> str:
        """Map emission type to component type."""
        mapping = {
            EmissionType.MODEL_USED: "model",
            EmissionType.TOOL_INVOKED: "tool",
            EmissionType.DATA_ACCESSED: "data_source",
            EmissionType.POLICY_APPLIED: "policy",
        }
        return mapping.get(emission_type, "tool")
