// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"
)

func TestResourcePolicyExemptionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *ResourcePolicyExemptionId
	}{
		{
			Name:  "empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "policy assignment in resource group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Authorization/policyExemptions/assignment1",
			Expected: &ResourcePolicyExemptionId{
				Name:       "assignment1",
				ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			},
		},
		{
			Name:  "policy assignment in resource group but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Authorization/policyExemptions/",
			Error: true,
		},
		{
			Name:  "policy assignment in subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyExemptions/assignment1",
			Expected: &ResourcePolicyExemptionId{
				Name:       "assignment1",
				ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "policy assignment in subscription but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyExemptions/",
			Error: true,
		},
		{
			Name:  "policy assignment in management group",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyExemptions/assignment1",
			Expected: &ResourcePolicyExemptionId{
				Name:       "assignment1",
				ResourceId: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "policy assignment in management group but no name",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyExemptions/",
			Error: true,
		},
		{
			Name:  "policy assignment in resource",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Authorization/policyExemptions/assignment1",
			Expected: &ResourcePolicyExemptionId{
				Name:       "assignment1",
				ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1",
			},
		},
		{
			Name:  "policy assignment in resource but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Authorization/policyExemptions/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ResourcePolicyExemptionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q", v.Expected.Name, actual.Name)
		}

		if v.Expected.ResourceId != actual.ResourceId {
			t.Fatalf("Expected %+v but got %+v", v.Expected.ResourceId, actual.ResourceId)
		}
	}
}
