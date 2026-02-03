package expander_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/consumer"
	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/expander"
	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/token"
)

type testCase struct {
	name        string
	targetDir   string
	inputFile   string
	outputFile  string
	shouldError bool
}

func TestFileExpander_ExpandToWriter(t *testing.T) {
	// Test cases
	testCases := []testCase{
		{
			name:        "single reference",
			targetDir:   "single-reference",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: false,
		},
		{
			name:        "recursive references",
			targetDir:   "recursive-references",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: false,
		},
		{
			name:        "url refer to not existing file",
			targetDir:   "url-refer-to-not-existing-file",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: true,
		},
		{
			name:        "cyclic reference",
			targetDir:   "cyclic-reference",
			inputFile:   "top.wrl",
			outputFile:  "merged.wrl",
			shouldError: true,
		},
		{
			name:        "reference to file in subdirectory",
			targetDir:   "reference-to-file-in-subdirectory",
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
	inputPath := filepath.Join("testdata", "input", tc.targetDir, tc.inputFile)
	expectedPath := filepath.Join("testdata", "expected", tc.targetDir, tc.outputFile)
	outputPath := filepath.Join("testdata", "output", tc.targetDir, tc.outputFile)

	// Check if input file exists
	if isNotFileExists(inputPath) {
		t.Fatalf("Input file not found: %s", inputPath)
	}

	// Create output directory if not exists
	prepareOutputDirectory(t, tc.targetDir)

	// Execute ExpandToWriter
	err := executeExpandToWriter(inputPath, outputPath)

	// Handle error cases
	if tc.shouldError {
		verifyExpectedError(t, err)
		return
	}

	// Handle success cases
	if err != nil {
		t.Fatalf("ExpandToWriter execution failed: %v", err)
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

func executeExpandToWriter(inputPath, outputPath string) error {
	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Create FileExpander with real implementations
	tokenizer := token.NewTokenizer()
	cons := consumer.NewInlineConsumer()
	fileExpander := expander.NewFileExpander(tokenizer, cons)

	// Expand the input file
	return fileExpander.ExpandToWriter(inputPath, outputFile)
}

func verifyExpectedError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Expected ExpandToWriter to fail, but it succeeded")
	}
	t.Logf("ExpandToWriter failed as expected: %v", err)
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
