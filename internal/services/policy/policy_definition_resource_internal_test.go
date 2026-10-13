// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestPolicyDefinitionParametersDiffRequiresReplacement(t *testing.T) {
	tests := []struct {
		name        string
		oldValue    string
		newValue    string
		requiresNew bool
		wantError   bool
	}{
		{name: "initial", newValue: `{"a":{"type":"String"}}`},
		{name: "unchanged", oldValue: `{"a":{"type":"String"}}`, newValue: `{"a":{"type":"String"}}`},
		{name: "added", oldValue: `{"a":{"type":"String"}}`, newValue: `{"a":{"type":"String"},"b":{"type":"String"}}`},
		{name: "removed", oldValue: `{"a":{"type":"String"},"b":{"type":"String"}}`, newValue: `{"a":{"type":"String"}}`, requiresNew: true},
		{name: "renamed", oldValue: `{"a":{"type":"String"}}`, newValue: `{"b":{"type":"String"}}`},
		{name: "removed while adding", oldValue: `{"a":{"type":"String"}}`, newValue: `{"b":{"type":"String"},"c":{"type":"String"}}`},
		{name: "cleared", oldValue: `{"a":{"type":"String"}}`, requiresNew: true},
		{name: "empty map cleared", oldValue: `{}`, requiresNew: true},
		{name: "invalid old JSON cleared", oldValue: `{`, requiresNew: true},
		{name: "invalid old JSON", oldValue: `{`, newValue: `{"a":{"type":"String"}}`, wantError: true},
		{name: "invalid new JSON", oldValue: `{"a":{"type":"String"}}`, newValue: `{`, wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resource := resourcePolicyDefinition()
			state := &terraform.InstanceState{
				ID: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/test",
				Attributes: map[string]string{
					"name":         "test",
					"policy_type":  "Custom",
					"mode":         "All",
					"display_name": "test",
					"parameters":   test.oldValue,
				},
			}
			config := terraform.NewResourceConfigRaw(map[string]any{
				"name":         "test",
				"policy_type":  "Custom",
				"mode":         "All",
				"display_name": "test",
				"parameters":   test.newValue,
			})
			diff, err := resource.Diff(context.Background(), state, config, nil)
			if test.wantError {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if requiresNew := diff != nil && diff.RequiresNew(); requiresNew != test.requiresNew {
				t.Fatalf("expected replacement=%t, got %t", test.requiresNew, requiresNew)
			}
		})
	}
}
