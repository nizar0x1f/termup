package ui

import (
	"testing"
)

func TestCleanPastedInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal text",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "URL with brackets",
			input:    "[https://example.com]",
			expected: "https://example.com",
		},
		{
			name:     "bracketed paste sequence",
			input:    "\x1b[200~https://example.com\x1b[201~",
			expected: "https://example.com",
		},
		{
			name:     "bracketed paste with visible brackets",
			input:    "\x1b[200~[https://example.com]\x1b[201~",
			expected: "https://example.com",
		},
		{
			name:     "text with control characters",
			input:    "hello\x00\x01world",
			expected: "helloworld",
		},
		{
			name:     "endpoint URL with brackets",
			input:    "[https://s3.amazonaws.com]",
			expected: "https://s3.amazonaws.com",
		},
		{
			name:     "access key with brackets",
			input:    "[AKIAIOSFODNN7EXAMPLE]",
			expected: "AKIAIOSFODNN7EXAMPLE",
		},
		{
			name:     "bucket name with brackets",
			input:    "[my-bucket-name]",
			expected: "my-bucket-name",
		},
		{
			name:     "short bucket name with brackets",
			input:    "[test]",
			expected: "test",
		},
		{
			name:     "bucket with underscores",
			input:    "[my_bucket_123]",
			expected: "my_bucket_123",
		},
		{
			name:     "brackets around very short text should be kept",
			input:    "[hi]",
			expected: "[hi]",
		},
		{
			name:     "brackets around non-bucket text should be kept",
			input:    "[hello!@#]",
			expected: "[hello!@#]",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "only brackets",
			input:    "[]",
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanPastedInput(tt.input)
			if result != tt.expected {
				t.Errorf("cleanPastedInput(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidBucketName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid bucket name",
			input:    "my-bucket",
			expected: true,
		},
		{
			name:     "valid bucket with numbers",
			input:    "bucket123",
			expected: true,
		},
		{
			name:     "valid bucket with underscores",
			input:    "my_bucket_name",
			expected: true,
		},
		{
			name:     "valid bucket with dots",
			input:    "my.bucket.name",
			expected: true,
		},
		{
			name:     "too short",
			input:    "ab",
			expected: false,
		},
		{
			name:     "invalid characters",
			input:    "bucket!@#",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidBucketName(tt.input)
			if result != tt.expected {
				t.Errorf("isValidBucketName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
