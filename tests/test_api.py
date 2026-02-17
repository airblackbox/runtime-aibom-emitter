"""Test FastAPI routes."""
import pytest
from fastapi.testclient import TestClient
from pkg.api.routes import router

client = TestClient(router)

def test_health():
    """Test health endpoint."""
    response = client.get("/v1/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "runtime-aibom-emitter"

def test_emit():
    """Test emit endpoint."""
    response = client.post(
        "/v1/emit",
        json={
            "emission_type": "model_used",
            "agent_id": "agent-1",
            "component_name": "GPT-4",
            "provider": "OpenAI"
        }
    )
    assert response.status_code == 200
    data = response.json()
    assert data["component_name"] == "GPT-4"

def test_emit_invalid_type():
    """Test emit with invalid type."""
    response = client.post(
        "/v1/emit",
        json={
            "emission_type": "invalid_type",
            "agent_id": "agent-1",
            "component_name": "Test"
        }
    )
    assert response.status_code == 400

@pytest.mark.asyncio
async def test_observe():
    """Test observe endpoint."""
    response = client.post(
        "/v1/observe",
        params={"agent_id": "agent-1", "limit": 5}
    )
    assert response.status_code == 200
    data = response.json()
    assert "emissions_generated" in data

def test_list_emissions():
    """Test listing emissions."""
    response = client.get("/v1/emissions")
    assert response.status_code == 200
    data = response.json()
    assert "count" in data
    assert "emissions" in data

@pytest.mark.asyncio
async def test_publish():
    """Test publish endpoint."""
    client.post(
        "/v1/emit",
        json={
            "emission_type": "model_used",
            "agent_id": "agent-1",
            "component_name": "GPT-4",
        }
    )
    response = client.post(
        "/v1/publish",
        params={"aibom_id": "aibom-1"}
    )
    assert response.status_code == 200
    data = response.json()
    assert data["published"] is True
