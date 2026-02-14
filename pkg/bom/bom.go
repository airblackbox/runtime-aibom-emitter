// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package bom defines the runtime AIBOM data model — an AI Bill of Materials
// that captures what models, tools, endpoints, and configurations were
// actually used at runtime, not just what was declared statically.
package bom

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AIBOM is an AI Bill of Materials combining static and runtime evidence.
type AIBOM struct {
	BOMFormat   string       `json:"bomFormat"`
	SpecVersion string       `json:"specVersion"`
	Version     int          `json:"version"`
	Metadata    Metadata     `json:"metadata"`
	Components  []Component  `json:"components"`
	Services    []Service    `json:"services,omitempty"`
	Evidence    Evidence     `json:"evidence"`
}

// Metadata describes the AIBOM itself.
type Metadata struct {
	Timestamp string `json:"timestamp"`
	ToolName  string `json:"tool_name"`
	ToolVer   string `json:"tool_version"`
	Component string `json:"component,omitempty"` // the app/agent being described
}

// Component is a model, library, or tool used by the AI system.
type Component struct {
	Type        string            `json:"type"` // model, library, tool, framework, endpoint
	Name        string            `json:"name"`
	Version     string            `json:"version,omitempty"`
	Provider    string            `json:"provider,omitempty"`
	BOMRef      string            `json:"bom-ref"`
	Properties  map[string]string `json:"properties,omitempty"`
}

// Service is an external endpoint called by the AI system.
type Service struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// Evidence records how the AIBOM was constructed.
type Evidence struct {
	Source     string          `json:"source"` // "runtime_traces", "static_scan", "hybrid"
	TraceIDs  []string        `json:"traceIds,omitempty"`
	SpanCount int             `json:"spanCount,omitempty"`
	Window    *TimeWindow     `json:"window,omitempty"`
	Findings  []RuntimeFinding `json:"findings,omitempty"`
}

// TimeWindow is the observation period for runtime evidence.
type TimeWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// RuntimeFinding is a specific observation from trace data.
type RuntimeFinding struct {
	Type       string `json:"type"` // model_usage, tool_call, endpoint_call
	Component  string `json:"component"`
	Count      int    `json:"count"`
	FirstSeen  string `json:"firstSeen"`
	LastSeen   string `json:"lastSeen"`
}

// New creates a new empty AIBOM.
func New(appName string) *AIBOM {
	return &AIBOM{
		BOMFormat:   "CycloneDX-AI",
		SpecVersion: "1.6",
		Version:     1,
		Metadata: Metadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			ToolName:  "runtime-aibom-emitter",
			ToolVer:   "0.1.0",
			Component: appName,
		},
		Evidence: Evidence{
			Source: "runtime_traces",
		},
	}
}

// AddModel adds a model component to the AIBOM.
func (b *AIBOM) AddModel(name, version, provider string) {
	ref := fmt.Sprintf("model-%s", bomRef(name, version))
	b.Components = append(b.Components, Component{
		Type:     "model",
		Name:     name,
		Version:  version,
		Provider: provider,
		BOMRef:   ref,
	})
}

// AddTool adds a tool component to the AIBOM.
func (b *AIBOM) AddTool(name, version string) {
	ref := fmt.Sprintf("tool-%s", bomRef(name, version))
	b.Components = append(b.Components, Component{
		Type:    "tool",
		Name:    name,
		Version: version,
		BOMRef:  ref,
	})
}

// AddFramework adds a framework component to the AIBOM.
func (b *AIBOM) AddFramework(name, version string) {
	ref := fmt.Sprintf("framework-%s", bomRef(name, version))
	b.Components = append(b.Components, Component{
		Type:    "framework",
		Name:    name,
		Version: version,
		BOMRef:  ref,
	})
}

// AddService adds an external service/endpoint.
func (b *AIBOM) AddService(name, endpoint, provider string) {
	b.Services = append(b.Services, Service{
		Name:     name,
		Endpoint: endpoint,
		Provider: provider,
	})
}

// AddFinding adds a runtime finding.
func (b *AIBOM) AddFinding(typ, component string, count int, firstSeen, lastSeen string) {
	b.Evidence.Findings = append(b.Evidence.Findings, RuntimeFinding{
		Type:      typ,
		Component: component,
		Count:     count,
		FirstSeen: firstSeen,
		LastSeen:  lastSeen,
	})
}

// ToJSON serializes the AIBOM to JSON.
func (b *AIBOM) ToJSON() ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}

// SaveFile writes the AIBOM to a JSON file.
func (b *AIBOM) SaveFile(path string) error {
	data, err := b.ToJSON()
	if err != nil {
		return fmt.Errorf("marshal AIBOM: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadFile reads an AIBOM from a JSON file.
func LoadFile(path string) (*AIBOM, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read AIBOM: %w", err)
	}
	var b AIBOM
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("parse AIBOM: %w", err)
	}
	return &b, nil
}

func bomRef(name, version string) string {
	h := sha256.Sum256([]byte(name + ":" + version))
	return fmt.Sprintf("%x", h[:8])
}
