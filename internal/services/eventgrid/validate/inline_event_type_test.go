package validate

import (
	"testing"
)

func TestValidateUniqueInlineEventTypeNames(t *testing.T) {
	testCases := []struct {
		name        string
		input       []interface{}
		shouldError bool
	}{
		{
			name:        "empty input",
			input:       []interface{}{},
			shouldError: false,
		},
		{
			name: "single event type",
			input: []interface{}{
				map[string]interface{}{
					"name": "EventType1",
				},
			},
			shouldError: false,
		},
		{
			name: "multiple unique event types",
			input: []interface{}{
				map[string]interface{}{
					"name": "EventType1",
				},
				map[string]interface{}{
					"name": "EventType2",
				},
				map[string]interface{}{
					"name": "EventType3",
				},
			},
			shouldError: false,
		},
		{
			name: "duplicate event type names",
			input: []interface{}{
				map[string]interface{}{
					"name": "EventType1",
				},
				map[string]interface{}{
					"name": "EventType2",
				},
				map[string]interface{}{
					"name": "EventType1",
				},
			},
			shouldError: true,
		},
		{
			name: "multiple duplicates",
			input: []interface{}{
				map[string]interface{}{
					"name": "EventType1",
				},
				map[string]interface{}{
					"name": "EventType1",
				},
				map[string]interface{}{
					"name": "EventType2",
				},
				map[string]interface{}{
					"name": "EventType2",
				},
			},
			shouldError: true,
		},
		{
			name: "case sensitive duplicate detection",
			input: []interface{}{
				map[string]interface{}{
					"name": "EventType1",
				},
				map[string]interface{}{
					"name": "eventtype1",
				},
			},
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUniqueInlineEventTypeNames(tc.input)

			if tc.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
