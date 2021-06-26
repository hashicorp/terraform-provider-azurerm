package validate

import (
	"strings"
	"testing"
)

func TestStorageQueueName_Validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "testing_123",
			ErrCount: 1,
		},
		{
			Value:    "testing123-",
			ErrCount: 1,
		},
		{
			Value:    "-testing123",
			ErrCount: 1,
		},
		{
			Value:    "TestingSG",
			ErrCount: 1,
		},
		{
			Value:    strings.Repeat("a", 256),
			ErrCount: 1,
		},
		{
			Value:    strings.Repeat("a", 1),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := StorageQueueName(tc.Value, "azurerm_storage_queue")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Storage Queue Name to trigger a validation error")
		}
	}
}
