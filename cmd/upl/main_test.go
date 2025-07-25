package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainFunction(t *testing.T) {
	// Skip in CI environments where TTY is not available
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping TTY-dependent test in CI environment")
	}

	tmpDir, err := os.MkdirTemp("", "upl-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath := filepath.Join(tmpDir, "upl-test")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "no arguments",
			args:           []string{},
			expectError:    true,
			expectedOutput: "Usage: upl <file-path>",
		},
		{
			name:           "non-existent file",
			args:           []string{"/path/to/non/existent/file.txt"},
			expectError:    true,
			expectedOutput: "no such file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpHome, err := os.MkdirTemp("", "home")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpHome)

			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "HOME="+tmpHome)

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none. Output: %s", outputStr)
			}

			if tt.expectedOutput != "" && !strings.Contains(outputStr, tt.expectedOutput) {

				if !strings.Contains(outputStr, "could not open a new TTY") {
					t.Errorf("Expected output to contain %q, got: %s", tt.expectedOutput, outputStr)
				}
			}
		})
	}
}

func TestMainWithValidFile(t *testing.T) {
	// Skip in CI environments where TTY is not available
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping TTY-dependent test in CI environment")
	}

	tmpFile, err := os.CreateTemp("", "test-upload.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("test content"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	tmpDir, err := os.MkdirTemp("", "upl-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath := filepath.Join(tmpDir, "upl-test")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tmpHome, err := os.MkdirTemp("", "home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	cmd = exec.Command(binaryPath, tmpFile.Name())
	cmd.Env = append(os.Environ(), "HOME="+tmpHome)

	output, _ := cmd.CombinedOutput()
	outputStr := string(output)

	if !strings.Contains(outputStr, "Enter Access Key ID") && !strings.Contains(outputStr, "could not open a new TTY") {
		t.Errorf("Expected to prompt for config or get TTY error, got: %s", outputStr)
	}
}
