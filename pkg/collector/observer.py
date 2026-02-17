"""Runtime observer for episode streams."""
from __future__ import annotations
import uuid
from typing import Any
from pkg.models.emission import Emission, EmissionType, EmissionSummary

class RuntimeObserver:
    """Observes episode execution and generates emissions."""
    def __init__(self, episode_store_url: str = "http://localhost:8000") -> None:
        self.episode_store_url = episode_store_url
        self._emissions: list[Emission] = []
        self._observed_episodes: set[str] = set()

    async def observe_episodes(
        self,
        agent_id: str,
        limit: int = 100
    ) -> list[Emission]:
        """Poll episodes and extract emissions."""
        new_emissions = []
        for i in range(limit):
            episode_data = self._simulate_episode(agent_id, i)
            if not episode_data:
                break
            emissions = self._extract_emissions(agent_id, episode_data)
            new_emissions.extend(emissions)
            self._emissions.extend(emissions)
        return new_emissions

    def _simulate_episode(
        self,
        agent_id: str,
        index: int
    ) -> dict[str, Any] | None:
        """Simulate fetching episode from store."""
        if index >= 5:
            return None
        episode_id = f"ep-{agent_id}-{index}"
        if episode_id in self._observed_episodes:
            return None
        self._observed_episodes.add(episode_id)
        return {
            "episode_id": episode_id,
            "agent_id": agent_id,
            "models_used": [
                {"name": "GPT-4", "version": "1.0", "provider": "OpenAI"}
            ],
            "tools_invoked": [
                {"name": "SearchTool", "version": "1.0", "provider": "Internal"}
            ],
            "data_accessed": [
                {"name": "UserDB", "version": "1.0", "provider": "Internal"}
            ],
            "timestamp": "2024-01-15T10:00:00Z",
        }

    def _extract_emissions(
        self,
        agent_id: str,
        episode_data: dict[str, Any]
    ) -> list[Emission]:
        """Extract emissions from episode data."""
        emissions = []
        for model in episode_data.get("models_used", []):
            emission = Emission(
                id=f"em-{uuid.uuid4().hex[:8]}",
                emission_type=EmissionType.MODEL_USED,
                agent_id=agent_id,
                component_name=model.get("name"),
                component_version=model.get("version", ""),
                provider=model.get("provider", ""),
                metadata={"episode_id": episode_data.get("episode_id")},
            )
            emissions.append(emission)
        for tool in episode_data.get("tools_invoked", []):
            emission = Emission(
                id=f"em-{uuid.uuid4().hex[:8]}",
                emission_type=EmissionType.TOOL_INVOKED,
                agent_id=agent_id,
                component_name=tool.get("name"),
                component_version=tool.get("version", ""),
                provider=tool.get("provider", ""),
                metadata={"episode_id": episode_data.get("episode_id")},
            )
            emissions.append(emission)
        for data in episode_data.get("data_accessed", []):
            emission = Emission(
                id=f"em-{uuid.uuid4().hex[:8]}",
                emission_type=EmissionType.DATA_ACCESSED,
                agent_id=agent_id,
                component_name=data.get("name"),
                component_version=data.get("version", ""),
                provider=data.get("provider", ""),
                metadata={"episode_id": episode_data.get("episode_id")},
            )
            emissions.append(emission)
        return emissions

    def get_summary(self, agent_id: str) -> EmissionSummary:
        """Get emission summary for agent."""
        agent_emissions = [
            e for e in self._emissions if e.agent_id == agent_id
        ]
        unique_models = set()
        unique_tools = set()
        unique_data = set()
        for emission in agent_emissions:
            if emission.emission_type == EmissionType.MODEL_USED:
                unique_models.add(emission.component_name)
            elif emission.emission_type == EmissionType.TOOL_INVOKED:
                unique_tools.add(emission.component_name)
            elif emission.emission_type == EmissionType.DATA_ACCESSED:
                unique_data.add(emission.component_name)
        return EmissionSummary(
            agent_id=agent_id,
            total_emissions=len(agent_emissions),
            unique_models=sorted(list(unique_models)),
            unique_tools=sorted(list(unique_tools)),
            unique_data_sources=sorted(list(unique_data)),
        )
