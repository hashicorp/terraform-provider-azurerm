package migration

import (
	"context"
	"reflect"
	"testing"
)

func TestSentinelAutomationRuleV0ToV1(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     `empty state`,
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: `state with "condition"`,
			input: map[string]interface{}{
				"condition": []interface{}{
					map[string]interface{}{
						"property": "FileName",
						"operator": "Contains",
						"values":   []interface{}{"foo"},
					},
				},
			},
			expected: map[string]interface{}{
				"condition_property": []interface{}{
					map[string]interface{}{
						"property": "FileName",
						"operator": "Contains",
						"values":   []interface{}{"foo"},
					},
				},
			},
		},
	}

	for _, test := range testData {
		t.Logf("Testing %q..", test.name)
		result, err := SentinelAutomationRuleV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
		if err != nil {
			t.Fatalf("Expected no error but got: %+v", err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Fatalf("expected %v but got %v!", test.expected, result)
		}
	}
}
