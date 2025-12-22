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

	// Test cases
	testCases := []struct {
		name        string
		sampleDir   string
		inputFile   string
		outputFile  string
		shouldError bool
	}{
		{
			name:        "sample1",
			sampleDir:   "sample1",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: false,
		},
		{
			name:        "sample2",
			sampleDir:   "sample2",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: false,
		},
		{
			name:        "sample3 - nonexistent referenced file",
			sampleDir:   "sample3",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: true,
		},
		{
			name:        "sample4 - cyclic reference",
			sampleDir:   "sample4",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare paths
			inputPath := filepath.Join("testdata", "input", tc.sampleDir, tc.inputFile)
			expectedPath := filepath.Join("testdata", "expected", tc.sampleDir, tc.outputFile)
			outputPath := filepath.Join("testdata", "output", tc.sampleDir, tc.outputFile)

			// Check if input file exists (skip if not)
			if _, err := os.Stat(inputPath); os.IsNotExist(err) {
				t.Skipf("Input file not found: %s", inputPath)
			}

			// Create output directory if not exists
			outputDir := filepath.Join("testdata", "output", tc.sampleDir)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				t.Fatalf("Failed to create output directory: %v", err)
			}

			// Execute CLI tool
			cmd := exec.Command("./vrml-inline-expander", inputPath, outputPath)
			output, err := cmd.CombinedOutput()

			if tc.shouldError {
				// Error is expected
				if err == nil {
					t.Fatalf("Expected CLI to fail, but it succeeded\nOutput: %s", output)
				}
				t.Logf("CLI failed as expected: %v\nOutput: %s", err, output)
				return
			}

			// No error expected
			if err != nil {
				t.Fatalf("CLI execution failed: %v\nOutput: %s", err, output)
			}

			// Check if expected file exists (skip if not)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Skipf("Expected file not found: %s", expectedPath)
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
				t.Errorf("Output does not match expected. See %s\n", outputPath)
			}
		})
	}
}
