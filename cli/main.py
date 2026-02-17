"""CLI commands for emissions."""
import click
from rich.console import Console
from rich.table import Table
import httpx
import json

console = Console()
BASE_URL = "http://localhost:8700/v1"

@click.group()
def cli():
    """Runtime AIBOM Emitter CLI."""
    pass

@cli.command()
def health():
    """Check emitter health."""
    try:
        with httpx.Client() as client:
            resp = client.get(f"{BASE_URL}/health")
            resp.raise_for_status()
            data = resp.json()
            console.print("[green]✓[/green] Emitter is healthy")
            console.print(f"  Emissions collected: {data['emissions_collected']}")
    except Exception as e:
        console.print(f"[red]✗[/red] Health check failed: {e}")

@cli.command()
@click.option("--agent-id", required=True, help="Agent ID")
@click.option("--limit", default=100, help="Episode limit")
def observe(agent_id: str, limit: int):
    """Observe agent episodes."""
    try:
        with httpx.Client() as client:
            resp = client.post(
                f"{BASE_URL}/observe",
                params={"agent_id": agent_id, "limit": limit}
            )
            resp.raise_for_status()
            data = resp.json()
            console.print(
                f"[green]✓[/green] Observed {data['emissions_generated']} "
                f"emissions from {agent_id}"
            )
            console.print(f"  Total collected: {data['total_collected']}")
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")

@cli.command()
@click.option("--agent-id", required=True, help="Agent ID")
def summary(agent_id: str):
    """Get agent emission summary."""
    try:
        with httpx.Client() as client:
            resp = client.get(f"{BASE_URL}/summary/{agent_id}")
            resp.raise_for_status()
            data = resp.json()
            table = Table(title=f"Emissions for {agent_id}")
            table.add_column("Metric")
            table.add_column("Value")
            table.add_row("Total Emissions", str(data["total_emissions"]))
            table.add_row("Unique Models", str(len(data["unique_models"])))
            table.add_row("Unique Tools", str(len(data["unique_tools"])))
            table.add_row("Unique Data Sources", str(len(data["unique_data_sources"])))
            console.print(table)
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")

@cli.command()
@click.option("--agent-id", default=None, help="Filter by agent")
@click.option("--type", "emission_type", default=None, help="Filter by type")
def list_emissions(agent_id: str, emission_type: str):
    """List collected emissions."""
    try:
        params = {}
        if agent_id:
            params["agent_id"] = agent_id
        if emission_type:
            params["emission_type"] = emission_type
        with httpx.Client() as client:
            resp = client.get(f"{BASE_URL}/emissions", params=params)
            resp.raise_for_status()
            data = resp.json()
            console.print(f"Found {data['count']} emissions")
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")

@cli.command()
@click.option("--aibom-id", required=True, help="AIBOM ID")
@click.option("--agent-id", default=None, help="Filter by agent")
def publish(aibom_id: str, agent_id: str):
    """Publish emissions to AIBOM."""
    try:
        params = {"aibom_id": aibom_id}
        if agent_id:
            params["agent_id"] = agent_id
        with httpx.Client() as client:
            resp = client.post(f"{BASE_URL}/publish", params=params)
            resp.raise_for_status()
            data = resp.json()
            console.print(f"[green]✓[/green] Published {data['count']} components")
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")

if __name__ == "__main__":
    cli()
