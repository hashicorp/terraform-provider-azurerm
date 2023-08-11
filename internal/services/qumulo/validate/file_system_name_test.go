package validate

import (
	"testing"
)

func TestFileSystemName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "a",
			ErrCount: 1,
		},
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "1---2",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "-helloWorld",
			ErrCount: 1,
		},
		{
			Value:    "helloWorld-",
			ErrCount: 1,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},
		{
			Value:    "123456789012345",
			ErrCount: 0,
		},
		{
			Value:    "1234567890123456",
			ErrCount: 1,
		},
		{
			Value:    "123456789012345-",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := FileSystemName(tc.Value, "azurerm_qumulo_file_system")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure Qumulo File System to trigger a validation error: %v", tc)
		}
	}
}
