// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"
)

func TestSecurityCenterAutomationV0ToV1(t *testing.T) {
	const canonical = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Security/automations/automation1"

	cases := []struct {
		input    string
		expected string
	}{
		{
			// already canonical IDs are left as-is
			input:    canonical,
			expected: canonical,
		},
		{
			// the legacy parser accepted an all-lowercase `resourcegroups` segment
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.Security/automations/automation1",
			expected: canonical,
		},
		{
			// the legacy parser did not validate the casing of the resource provider
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/microsoft.security/automations/automation1",
			expected: canonical,
		},
	}

	for _, tc := range cases {
		rawState := map[string]interface{}{
			"id": tc.input,
		}

		result, err := SecurityCenterAutomationV0ToV1{}.UpgradeFunc()(context.Background(), rawState, nil)
		if err != nil {
			t.Fatalf("upgrading %q: %+v", tc.input, err)
		}

		if result["id"] != tc.expected {
			t.Errorf("expected %q to migrate to %q, got %q", tc.input, tc.expected, result["id"])
		}
	}
}
