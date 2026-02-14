// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

package bom

import (
	"encoding/json"
	"testing"
)

func TestNewAIBOM(t *testing.T) {
	b := New("test-agent")
	if b.BOMFormat != "CycloneDX-AI" {
		t.Fatalf("expected CycloneDX-AI, got %s", b.BOMFormat)
	}
	if b.Metadata.Component != "test-agent" {
		t.Fatalf("expected test-agent, got %s", b.Metadata.Component)
	}
	if b.Evidence.Source != "runtime_traces" {
		t.Fatalf("expected runtime_traces, got %s", b.Evidence.Source)
	}
}

func TestAddModel(t *testing.T) {
	b := New("test")
	b.AddModel("gpt-4", "2024-01-01", "openai")
	if len(b.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(b.Components))
	}
	c := b.Components[0]
	if c.Type != "model" || c.Name != "gpt-4" || c.Provider != "openai" {
		t.Fatalf("unexpected component: %+v", c)
	}
}

func TestAddTool(t *testing.T) {
	b := New("test")
	b.AddTool("web_search", "1.0")
	if len(b.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(b.Components))
	}
	if b.Components[0].Type != "tool" {
		t.Fatalf("expected tool type, got %s", b.Components[0].Type)
	}
}

func TestAddService(t *testing.T) {
	b := New("test")
	b.AddService("openai-api", "https://api.openai.com", "openai")
	if len(b.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(b.Services))
	}
}

func TestAddFinding(t *testing.T) {
	b := New("test")
	b.AddFinding("model_usage", "gpt-4", 5, "2024-01-01T00:00:00Z", "2024-01-01T01:00:00Z")
	if len(b.Evidence.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(b.Evidence.Findings))
	}
	f := b.Evidence.Findings[0]
	if f.Type != "model_usage" || f.Component != "gpt-4" || f.Count != 5 {
		t.Fatalf("unexpected finding: %+v", f)
	}
}

func TestToJSON(t *testing.T) {
	b := New("test")
	b.AddModel("gpt-4", "", "openai")
	b.AddTool("web_search", "")
	data, err := b.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["bomFormat"] != "CycloneDX-AI" {
		t.Fatal("missing bomFormat")
	}
}

func TestCompleteAIBOM(t *testing.T) {
	b := New("my-agent")
	b.AddModel("gpt-4", "2024-01", "openai")
	b.AddModel("claude-3-sonnet", "2024-02", "anthropic")
	b.AddTool("web_search", "")
	b.AddTool("code_interpreter", "")
	b.AddFramework("langchain", "0.1.0")
	b.AddService("openai", "https://api.openai.com/v1", "openai")
	b.AddFinding("model_usage", "gpt-4", 10, "2024-01-01T00:00:00Z", "2024-01-01T02:00:00Z")
	b.Evidence.SpanCount = 50
	b.Evidence.TraceIDs = []string{"abc123", "def456"}

	if len(b.Components) != 5 {
		t.Fatalf("expected 5 components, got %d", len(b.Components))
	}
	if len(b.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(b.Services))
	}

	data, err := b.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("empty JSON")
	}
}
