// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package collector extracts AIBOM components from OTLP JSON trace exports.
// It scans span attributes for GenAI semantic convention fields and builds
// a runtime inventory of models, tools, and endpoints actually used.
package collector

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nostalgicskinco/runtime-aibom-emitter/pkg/bom"
)

// Config controls what the collector extracts from traces.
type Config struct {
	ModelAttrs    []string `json:"model_attrs"`
	ProviderAttrs []string `json:"provider_attrs"`
	ToolAttrs     []string `json:"tool_attrs"`
	EndpointAttrs []string `json:"endpoint_attrs"`
}

// DefaultConfig returns a config targeting standard GenAI semantic conventions.
func DefaultConfig() Config {
	return Config{
		ModelAttrs:    []string{"gen_ai.request.model", "llm.model_name", "ai.model.id"},
		ProviderAttrs: []string{"gen_ai.system", "llm.provider", "ai.provider"},
		ToolAttrs:     []string{"tool.name", "mcp.tool.name", "gen_ai.tool.name"},
		EndpointAttrs: []string{"server.address", "url.full", "http.url"},
	}
}

// otlpExport mirrors OTLP JSON.
type otlpExport struct {
	ResourceSpans []resourceSpan `json:"resourceSpans"`
}
type resourceSpan struct {
	Resource   *resource   `json:"resource,omitempty"`
	ScopeSpans []scopeSpan `json:"scopeSpans"`
}
type resource struct {
	Attributes []kv `json:"attributes,omitempty"`
}
type scopeSpan struct {
	Scope *scope    `json:"scope,omitempty"`
	Spans []rawSpan `json:"spans"`
}
type scope struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}
type rawSpan struct {
	TraceID    string `json:"traceId"`
	Name       string `json:"name"`
	StartNano  string `json:"startTimeUnixNano"`
	Attributes []kv   `json:"attributes,omitempty"`
}
type kv struct {
	Key   string  `json:"key"`
	Value kvValue `json:"value"`
}
type kvValue struct {
	StringValue *string `json:"stringValue,omitempty"`
}

func getStr(attrs []kv, candidates []string) string {
	for _, c := range candidates {
		for _, a := range attrs {
			if a.Key == c && a.Value.StringValue != nil {
				return *a.Value.StringValue
			}
		}
	}
	return ""
}

type componentKey struct {
	typ  string
	name string
}

// CollectFromFile reads an OTLP JSON file and produces an AIBOM.
func CollectFromFile(path string, appName string, cfg Config) (*bom.AIBOM, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read trace file: %w", err)
	}
	return CollectFromJSON(data, appName, cfg)
}

// CollectFromJSON parses OTLP JSON and produces an AIBOM.
func CollectFromJSON(data []byte, appName string, cfg Config) (*bom.AIBOM, error) {
	var export otlpExport
	if err := json.Unmarshal(data, &export); err != nil {
		return nil, fmt.Errorf("parse OTLP: %w", err)
	}

	b := bom.New(appName)
	seen := make(map[componentKey]bool)
	traceIDs := make(map[string]bool)
	spanCount := 0
	var earliest, latest time.Time

	for _, rs := range export.ResourceSpans {
		// Extract framework from scope
		for _, ss := range rs.ScopeSpans {
			if ss.Scope != nil && ss.Scope.Name != "" {
				fwKey := componentKey{typ: "framework", name: ss.Scope.Name}
				if !seen[fwKey] {
					seen[fwKey] = true
					b.AddFramework(ss.Scope.Name, ss.Scope.Version)
				}
			}

			for _, span := range ss.Spans {
				spanCount++
				traceIDs[span.TraceID] = true

				// Parse timestamp
				ts := parseNanos(span.StartNano)
				if !ts.IsZero() {
					if earliest.IsZero() || ts.Before(earliest) {
						earliest = ts
					}
					if latest.IsZero() || ts.After(latest) {
						latest = ts
					}
				}

				// Extract model
				model := getStr(span.Attributes, cfg.ModelAttrs)
				provider := getStr(span.Attributes, cfg.ProviderAttrs)
				if model != "" {
					mKey := componentKey{typ: "model", name: model}
					if !seen[mKey] {
						seen[mKey] = true
						b.AddModel(model, "", provider)
						b.AddFinding("model_usage", model, 1, ts.Format(time.RFC3339), ts.Format(time.RFC3339))
					}
				}

				// Extract tool
				tool := getStr(span.Attributes, cfg.ToolAttrs)
				if tool != "" {
					tKey := componentKey{typ: "tool", name: tool}
					if !seen[tKey] {
						seen[tKey] = true
						b.AddTool(tool, "")
						b.AddFinding("tool_call", tool, 1, ts.Format(time.RFC3339), ts.Format(time.RFC3339))
					}
				}

				// Extract endpoint
				endpoint := getStr(span.Attributes, cfg.EndpointAttrs)
				if endpoint != "" {
					// Normalize endpoint to host
					host := endpoint
					if idx := strings.Index(endpoint, "://"); idx >= 0 {
						host = endpoint[idx+3:]
					}
					if idx := strings.Index(host, "/"); idx >= 0 {
						host = host[:idx]
					}
					eKey := componentKey{typ: "endpoint", name: host}
					if !seen[eKey] {
						seen[eKey] = true
						b.AddService(host, endpoint, provider)
					}
				}
			}
		}
	}

	// Set evidence metadata
	for tid := range traceIDs {
		b.Evidence.TraceIDs = append(b.Evidence.TraceIDs, tid)
	}
	b.Evidence.SpanCount = spanCount
	if !earliest.IsZero() && !latest.IsZero() {
		b.Evidence.Window = &bom.TimeWindow{
			Start: earliest.Format(time.RFC3339),
			End:   latest.Format(time.RFC3339),
		}
	}

	return b, nil
}

func parseNanos(s string) time.Time {
	var ns int64
	fmt.Sscanf(s, "%d", &ns)
	if ns == 0 {
		return time.Time{}
	}
	return time.Unix(0, ns)
}
