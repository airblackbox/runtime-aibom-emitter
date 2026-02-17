"""FastAPI routes for emissions."""
from __future__ import annotations
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from pkg.models.emission import Emission, EmissionType, EmissionSummary
from pkg.collector.observer import RuntimeObserver
from pkg.emitter.publisher import EmissionPublisher

router = FastAPI(title="Runtime AIBOM Emitter")
observer = RuntimeObserver()
publisher = EmissionPublisher()

class EmissionInput(BaseModel):
    """Emission input model."""
    emission_type: str
    agent_id: str
    component_name: str
    component_version: str = ""
    provider: str = ""

@router.get("/v1/health")
async def health():
    """Health check endpoint."""
    return {
        "status": "ok",
        "service": "runtime-aibom-emitter",
        "emissions_collected": len(publisher._emissions)
    }

@router.post("/v1/emit")
async def emit(input_data: EmissionInput) -> Emission:
    """Record an emission event."""
    try:
        emission_type = EmissionType[input_data.emission_type.upper()]
    except KeyError:
        raise HTTPException(
            status_code=400,
            detail=f"Invalid emission type: {input_data.emission_type}"
        )
    emission = Emission(
        emission_type=emission_type,
        agent_id=input_data.agent_id,
        component_name=input_data.component_name,
        component_version=input_data.component_version,
        provider=input_data.provider,
    )
    publisher.collect(emission)
    return emission

@router.post("/v1/observe")
async def observe(agent_id: str, limit: int = 100):
    """Observe episodes and generate emissions."""
    emissions = await observer.observe_episodes(agent_id, limit)
    publisher.collect_batch(emissions)
    return {
        "agent_id": agent_id,
        "emissions_generated": len(emissions),
        "total_collected": len(publisher._emissions)
    }

@router.get("/v1/summary/{agent_id}")
async def get_summary(agent_id: str) -> EmissionSummary:
    """Get emission summary for agent."""
    summary = observer.get_summary(agent_id)
    return summary

@router.get("/v1/emissions")
async def list_emissions(
    agent_id: str | None = None,
    emission_type: str | None = None
):
    """List collected emissions."""
    type_filter = None
    if emission_type:
        try:
            type_filter = EmissionType[emission_type.upper()]
        except KeyError:
            pass
    emissions = publisher.get_emissions(agent_id, type_filter)
    return {
        "count": len(emissions),
        "emissions": emissions
    }

@router.post("/v1/publish")
async def publish_to_aibom(aibom_id: str, agent_id: str | None = None):
    """Publish emissions to AIBOM engine."""
    result = await publisher.publish(aibom_id, agent_id)
    return result

@router.post("/v1/export")
async def export_emissions(filepath: str = "/tmp/emissions.json"):
    """Export emissions to JSON file."""
    publisher.export_json(filepath)
    return {
        "exported": True,
        "filepath": filepath,
        "count": len(publisher._emissions)
    }
