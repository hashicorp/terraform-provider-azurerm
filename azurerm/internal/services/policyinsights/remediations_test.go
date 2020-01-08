package policyinsights

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestParseScope(t *testing.T) {
	testData := []struct {
		input    string
		expected *RemediationScope
	}{
		{
			// resource group as scope
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			expected: &RemediationScope{
				Type:              AtResourceGroup,
				Scope:             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
				SubscriptionId:    utils.String("00000000-0000-0000-0000-000000000000"),
				ManagementGroupId: nil,
				ResourceGroup:     utils.String("foo"),
			},
		},
		{
			// resource group as scope but no name
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			expected: nil,
		},
		{
			// subscription as scope
			input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			expected: &RemediationScope{
				Type:              AtSubscription,
				Scope:             "/subscriptions/00000000-0000-0000-0000-000000000000",
				SubscriptionId:    utils.String("00000000-0000-0000-0000-000000000000"),
				ManagementGroupId: nil,
				ResourceGroup:     nil,
			},
		},
		{
			// management group id as scope
			input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			expected: &RemediationScope{
				Type:              AtManagementGroup,
				Scope:             "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
				SubscriptionId:    nil,
				ManagementGroupId: utils.String("00000000-0000-0000-0000-000000000000"),
				ResourceGroup:     nil,
			},
		},
		{
			// management group id as scope with different casing
			input: "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000",
			expected: &RemediationScope{
				Type:              AtManagementGroup,
				Scope:             "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000",
				SubscriptionId:    nil,
				ManagementGroupId: utils.String("00000000-0000-0000-0000-000000000000"),
				ResourceGroup:     nil,
			},
		},
		{
			// resource id as scope
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1",
			expected: &RemediationScope{
				Type:              AtResource,
				Scope:             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1",
				SubscriptionId:    utils.String("00000000-0000-0000-0000-000000000000"),
				ManagementGroupId: nil,
				ResourceGroup:     utils.String("foo"),
			},
		},
		{
			// illegal resource id
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines",
			expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		actual, err := ParseScope(v.input)
		if err != nil {
			if v.expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if !reflect.DeepEqual(v.expected, actual) {
			t.Fatalf("Expected %v, but got %v", *v.expected, *actual)
		}
	}
}

func TestParseRemediationId(t *testing.T) {
	testData := []struct {
		input    string
		expected *RemediationId
	}{
		{
			// resource group as scope
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/microsoft.policyinsights/remediations/remediation1",
			expected: &RemediationId{
				Name: "remediation1",
				RemediationScope: RemediationScope{
					Type:              AtResourceGroup,
					Scope:             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
					SubscriptionId:    utils.String("00000000-0000-0000-0000-000000000000"),
					ManagementGroupId: nil,
					ResourceGroup:     utils.String("foo"),
				},
			},
		},
		{
			// resource group as scope but no name
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/microsoft.policyinsights/remediations/",
			expected: nil,
		},
		{
			// subscription as scope
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/microsoft.policyinsights/remediations/remediation1",
			expected: &RemediationId{
				Name: "remediation1",
				RemediationScope: RemediationScope{
					Type:              AtSubscription,
					Scope:             "/subscriptions/00000000-0000-0000-0000-000000000000",
					SubscriptionId:    utils.String("00000000-0000-0000-0000-000000000000"),
					ManagementGroupId: nil,
					ResourceGroup:     nil,
				},
			},
		},
		{
			// management group id as scope
			input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/microsoft.policyinsights/remediations/remediation1",
			expected: &RemediationId{
				Name: "remediation1",
				RemediationScope: RemediationScope{
					Type:              AtManagementGroup,
					Scope:             "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
					SubscriptionId:    nil,
					ManagementGroupId: utils.String("00000000-0000-0000-0000-000000000000"),
					ResourceGroup:     nil,
				},
			},
		},
		{
			// management group id as scope with different casing
			input: "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/microsoft.policyinsights/remediations/remediation1",
			expected: &RemediationId{
				Name: "remediation1",
				RemediationScope: RemediationScope{
					Type:              AtManagementGroup,
					Scope:             "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000",
					SubscriptionId:    nil,
					ManagementGroupId: utils.String("00000000-0000-0000-0000-000000000000"),
					ResourceGroup:     nil,
				},
			},
		},
		{
			// resource id as scope
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1/providers/microsoft.policyinsights/remediations/remediation1",
			expected: &RemediationId{
				Name: "remediation1",
				RemediationScope: RemediationScope{
					Type:              AtResource,
					Scope:             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1",
					SubscriptionId:    utils.String("00000000-0000-0000-0000-000000000000"),
					ManagementGroupId: nil,
					ResourceGroup:     utils.String("foo"),
				},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		actual, err := ParseRemediationId(v.input)
		if err != nil {
			if v.expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if !reflect.DeepEqual(v.expected, actual) {
			t.Fatalf("Expected %v, but got %v", *v.expected, *actual)
		}
	}
}
