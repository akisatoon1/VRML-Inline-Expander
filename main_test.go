package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type testCase struct {
	name        string
	sampleDir   string
	inputFile   string
	outputFile  string
	shouldError bool
}

func TestCLI(t *testing.T) {
	// Build the CLI tool
	buildCmd := exec.Command("go", "build", "-o", "vrml-inline-expander")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI tool: %v", err)
	}
	defer os.Remove("vrml-inline-expander")

	// Test cases
	testCases := []testCase{
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
		{
			name:        "sample5",
			sampleDir:   "sample5",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}

func runTestCase(t *testing.T, tc testCase) {
	// Prepare paths
	inputPath := filepath.Join("testdata", "input", tc.sampleDir, tc.inputFile)
	expectedPath := filepath.Join("testdata", "expected", tc.sampleDir, tc.outputFile)
	outputPath := filepath.Join("testdata", "output", tc.sampleDir, tc.outputFile)

	// Check if input file exists
	if isNotFileExists(inputPath) {
		t.Fatalf("Input file not found: %s", inputPath)
	}

	// Create output directory if not exists
	prepareOutputDirectory(t, tc.sampleDir)

	// Execute CLI tool
	output, err := executeCLI(inputPath, outputPath)

	// Handle error cases
	if tc.shouldError {
		verifyExpectedError(t, output, err)
		return
	}

	// Handle success cases
	if err != nil {
		t.Fatalf("CLI execution failed: %v\nOutput: %s", err, output)
	}

	compareOutputFileWithExpected(t, outputPath, expectedPath)
}

func isNotFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}

func prepareOutputDirectory(t *testing.T, sampleDir string) {
	outputDir := filepath.Join("testdata", "output", sampleDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
}

func executeCLI(inputPath, outputPath string) ([]byte, error) {
	cmd := exec.Command("./vrml-inline-expander", inputPath, outputPath)
	return cmd.CombinedOutput()
}

func verifyExpectedError(t *testing.T, output []byte, err error) {
	if err == nil {
		t.Fatalf("Expected CLI to fail, but it succeeded\nOutput: %s", output)
	}
	t.Logf("CLI failed as expected: %v\nOutput: %s", err, output)
}

func compareOutputFileWithExpected(t *testing.T, outputPath, expectedPath string) {
	// Check if expected file exists
	if isNotFileExists(expectedPath) {
		t.Fatalf("Expected file not found: %s", expectedPath)
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
}
