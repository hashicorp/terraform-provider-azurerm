package validate

import (
	"strings"
	"testing"
)

func TestStorageBlobIndexTagName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		strings.Repeat("w", 128),
	}
	for _, v := range validNames {
		_, errors := StorageBlobIndexTagName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Blob Index Tag Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"",
		strings.Repeat("w", 129),
	}
	for _, v := range invalidNames {
		if _, errors := StorageBlobIndexTagName(v, "name"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid Blob Index Tag Name", v)
		}
	}
}

func TestStorageBlobIndexTagValue(t *testing.T) {
	validNames := []string{
		"valid-name",
		"",
		strings.Repeat("w", 256),
	}
	for _, v := range validNames {
		_, errors := StorageBlobIndexTagValue(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Blob Index Tag Value: %q", v, errors)
		}
	}

	invalidNames := []string{
		strings.Repeat("w", 257),
	}
	for _, v := range invalidNames {
		if _, errors := StorageBlobIndexTagValue(v, "name"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid Blob Index Tag Value", v)
		}
	}
}
