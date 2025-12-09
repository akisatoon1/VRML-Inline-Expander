package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLI(t *testing.T) {
	// Build the CLI tool
	buildCmd := exec.Command("go", "build", "-o", "vrml-inline-expander")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI tool: %v", err)
	}
	defer os.Remove("vrml-inline-expander")

	// Prepare paths
	inputPath := filepath.Join("testdata", "input", "sample.wrl")
	expectedPath := filepath.Join("testdata", "expected", "sample.wrl")
	outputPath := filepath.Join("testdata", "output", "sample.wrl")

	// Create output directory if not exists
	outputDir := filepath.Join("testdata", "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Execute CLI tool
	cmd := exec.Command("./vrml-inline-expander", inputPath, outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI execution failed: %v\nOutput: %s", err, output)
	}

	// Read output file
	outputContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Read expected file
	expectedContent, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read expected file: %v", err)
	}

	// Compare
	if string(outputContent) != string(expectedContent) {
		t.Errorf("Output does not match expected. See testdata/output/\n")
	}
}
