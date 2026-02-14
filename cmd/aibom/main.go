// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command aibom generates a runtime AIBOM from OTLP JSON trace exports.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nostalgicskinco/runtime-aibom-emitter/pkg/collector"
)

func main() {
	traceFile := flag.String("trace", "", "Path to OTLP JSON trace export")
	appName := flag.String("app", "unknown", "Application/agent name")
	output := flag.String("output", "", "Output file (default: stdout)")
	flag.Parse()

	if *traceFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: aibom -trace <traces.json> [-app <name>] [-output <file>]\n")
		os.Exit(1)
	}

	cfg := collector.DefaultConfig()
	b, err := collector.CollectFromFile(*traceFile, *appName, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		if err := b.SaveFile(*output); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "AIBOM written to %s (%d components, %d services)\n", *output, len(b.Components), len(b.Services))
	} else {
		data, err := b.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(data)
	}
}
