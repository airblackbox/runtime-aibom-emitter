// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

package collector

import (
	"testing"
)

const sampleTrace = `{
	"resourceSpans": [{
		"scopeSpans": [{
			"scope": {"name": "langchain", "version": "0.1.0"},
			"spans": [
				{
					"traceId": "abc123",
					"name": "chat gpt-4",
					"startTimeUnixNano": "1700000000000000000",
					"attributes": [
						{"key": "gen_ai.request.model", "value": {"stringValue": "gpt-4"}},
						{"key": "gen_ai.system", "value": {"stringValue": "openai"}},
						{"key": "server.address", "value": {"stringValue": "api.openai.com"}}
					]
				},
				{
					"traceId": "abc123",
					"name": "tool web_search",
					"startTimeUnixNano": "1700000000100000000",
					"attributes": [
						{"key": "tool.name", "value": {"stringValue": "web_search"}}
					]
				},
				{
					"traceId": "abc123",
					"name": "chat claude-3-sonnet",
					"startTimeUnixNano": "1700000000200000000",
					"attributes": [
						{"key": "gen_ai.request.model", "value": {"stringValue": "claude-3-sonnet"}},
						{"key": "gen_ai.system", "value": {"stringValue": "anthropic"}},
						{"key": "server.address", "value": {"stringValue": "api.anthropic.com"}}
					]
				}
			]
		}]
	}]
}`

func TestCollectFromJSON(t *testing.T) {
	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(sampleTrace), "test-agent", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}

	if b.Metadata.Component != "test-agent" {
		t.Fatalf("expected test-agent, got %s", b.Metadata.Component)
	}
}

func TestExtractsModels(t *testing.T) {
	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(sampleTrace), "test", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}

	models := 0
	for _, c := range b.Components {
		if c.Type == "model" {
			models++
		}
	}
	if models != 2 {
		t.Fatalf("expected 2 models, got %d", models)
	}
}

func TestExtractsTools(t *testing.T) {
	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(sampleTrace), "test", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}

	tools := 0
	for _, c := range b.Components {
		if c.Type == "tool" {
			tools++
		}
	}
	if tools != 1 {
		t.Fatalf("expected 1 tool, got %d", tools)
	}
}

func TestExtractsFramework(t *testing.T) {
	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(sampleTrace), "test", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}

	frameworks := 0
	for _, c := range b.Components {
		if c.Type == "framework" {
			frameworks++
		}
	}
	if frameworks != 1 {
		t.Fatalf("expected 1 framework, got %d", frameworks)
	}
}

func TestExtractsServices(t *testing.T) {
	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(sampleTrace), "test", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}

	if len(b.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(b.Services))
	}
}

func TestEvidenceMetadata(t *testing.T) {
	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(sampleTrace), "test", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}

	if b.Evidence.SpanCount != 3 {
		t.Fatalf("expected 3 spans, got %d", b.Evidence.SpanCount)
	}
	if len(b.Evidence.TraceIDs) != 1 {
		t.Fatalf("expected 1 trace ID, got %d", len(b.Evidence.TraceIDs))
	}
	if b.Evidence.Window == nil {
		t.Fatal("expected time window")
	}
}

func TestDeduplicatesModels(t *testing.T) {
	trace := `{
		"resourceSpans": [{
			"scopeSpans": [{
				"spans": [
					{
						"traceId": "t1",
						"name": "chat gpt-4",
						"startTimeUnixNano": "1700000000000000000",
						"attributes": [{"key": "gen_ai.request.model", "value": {"stringValue": "gpt-4"}}]
					},
					{
						"traceId": "t1",
						"name": "chat gpt-4 again",
						"startTimeUnixNano": "1700000000100000000",
						"attributes": [{"key": "gen_ai.request.model", "value": {"stringValue": "gpt-4"}}]
					}
				]
			}]
		}]
	}`

	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(trace), "test", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}

	models := 0
	for _, c := range b.Components {
		if c.Type == "model" {
			models++
		}
	}
	if models != 1 {
		t.Fatalf("expected 1 deduplicated model, got %d", models)
	}
}

func TestEmptyTrace(t *testing.T) {
	trace := `{"resourceSpans": []}`
	cfg := DefaultConfig()
	b, err := CollectFromJSON([]byte(trace), "test", cfg)
	if err != nil {
		t.Fatalf("CollectFromJSON: %v", err)
	}
	if len(b.Components) != 0 {
		t.Fatalf("expected 0 components, got %d", len(b.Components))
	}
}
