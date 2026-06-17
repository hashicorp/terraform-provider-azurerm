// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	batchRules "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
)

func TestBatchRuleSetOriginGroupOverridesMatchDesired(t *testing.T) {
	tests := []struct {
		name     string
		actual   azuresdkhacks.BatchRuleSetResource
		desired  azuresdkhacks.BatchRuleSetResource
		expected bool
	}{
		{
			name:     "origin group override removed not yet reflected",
			actual:   batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Cdn/profiles/profile/originGroups/group-a")),
			desired:  batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "")),
			expected: false,
		},
		{
			name:     "origin group override removed reflected in readback",
			actual:   batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "")),
			desired:  batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "")),
			expected: true,
		},
		{
			name:     "origin group override add reflected in readback",
			actual:   batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Cdn/profiles/profile/originGroups/group-a")),
			desired:  batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Cdn/profiles/profile/originGroups/group-a")),
			expected: true,
		},
		{
			name:     "origin group override points at wrong target",
			actual:   batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Cdn/profiles/profile/originGroups/group-b")),
			desired:  batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Cdn/profiles/profile/originGroups/group-a")),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := batchRuleSetOriginGroupOverridesMatchDesired(&test.actual, test.desired)
			if err != nil {
				t.Fatalf("expected no error but got %q", err)
			}

			if actual != test.expected {
				t.Fatalf("expected %t but got %t", test.expected, actual)
			}
		})
	}
}

func TestBatchRuleSetStatusesSettled(t *testing.T) {
	tests := []struct {
		name        string
		input       azuresdkhacks.BatchRuleSetResource
		expected    bool
		expectError bool
	}{
		{
			name:     "missing top level statuses does not block settled rule",
			input:    batchRuleSetResourceForTest(batchRuleSetRuleForTest("rule-a", "")),
			expected: false,
		},
		{
			name:     "top level in progress blocks",
			input:    batchRuleSetResourceWithStatusesForTest(pointer.To("Creating"), pointer.To("Deploying"), batchRuleSetRuleForTest("rule-a", "")),
			expected: false,
		},
		{
			name:     "top level deployment status not started does not block once provisioning succeeds",
			input:    batchRuleSetResourceWithStatusesForTest(pointer.To(string(batchRules.AfdProvisioningStateSucceeded)), pointer.To("NotStarted"), batchRuleSetRuleForTest("rule-a", "")),
			expected: true,
		},
		{
			name:        "top level failed errors",
			input:       batchRuleSetResourceWithStatusesForTest(pointer.To(string(batchRules.AfdProvisioningStateFailed)), pointer.To(string(batchRules.DeploymentStatusFailed)), batchRuleSetRuleForTest("rule-a", "")),
			expectError: true,
		},
		{
			name:     "rule in progress blocks",
			input:    batchRuleSetResourceForTest(batchRuleSetRuleWithStatusesForTest("rule-a", "", pointer.To(batchRules.AfdProvisioningStateUpdating), pointer.To(batchRules.DeploymentStatusInProgress))),
			expected: false,
		},
		{
			name:        "rule failed errors",
			input:       batchRuleSetResourceWithStatusesForTest(pointer.To(string(batchRules.AfdProvisioningStateSucceeded)), nil, batchRuleSetRuleWithStatusesForTest("rule-a", "", pointer.To(batchRules.AfdProvisioningStateFailed), pointer.To(batchRules.DeploymentStatusFailed))),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := batchRuleSetStatusesSettled(&test.input)
			if test.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error but got %q", err)
			}

			if actual != test.expected {
				t.Fatalf("expected %t but got %t", test.expected, actual)
			}
		})
	}
}

func batchRuleSetResourceForTest(rules ...azuresdkhacks.BatchRuleProperties) azuresdkhacks.BatchRuleSetResource {
	return azuresdkhacks.BatchRuleSetResource{
		Properties: &azuresdkhacks.BatchRuleSetProperties{
			Rules: &rules,
		},
	}
}

func batchRuleSetResourceWithStatusesForTest(provisioningState *string, deploymentStatus *string, rules ...azuresdkhacks.BatchRuleProperties) azuresdkhacks.BatchRuleSetResource {
	return azuresdkhacks.BatchRuleSetResource{
		Properties: &azuresdkhacks.BatchRuleSetProperties{
			ProvisioningState: provisioningState,
			DeploymentStatus:  deploymentStatus,
			Rules:             &rules,
		},
	}
}

func batchRuleSetRuleForTest(name string, originGroupID string) azuresdkhacks.BatchRuleProperties {
	return batchRuleSetRuleWithStatusesForTest(name, originGroupID, nil, nil)
}

func batchRuleSetRuleWithStatusesForTest(name string, originGroupID string, provisioningState *batchRules.AfdProvisioningState, deploymentStatus *batchRules.DeploymentStatus) azuresdkhacks.BatchRuleProperties {
	actions := []batchRules.DeliveryRuleAction{
		batchRules.DeliveryRuleRouteConfigurationOverrideAction{
			Name: batchRules.DeliveryRuleActionNameRouteConfigurationOverride,
			Parameters: batchRules.RouteConfigurationOverrideActionParameters{
				TypeName:            batchRules.DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters,
				OriginGroupOverride: originGroupOverrideForTest(originGroupID),
			},
		},
	}

	return azuresdkhacks.BatchRuleProperties{
		Name:              pointer.To(name),
		RuleName:          pointer.To(name),
		Actions:           &actions,
		ProvisioningState: provisioningState,
		DeploymentStatus:  deploymentStatus,
	}
}

func originGroupOverrideForTest(id string) *batchRules.OriginGroupOverride {
	if id == "" {
		return nil
	}

	return &batchRules.OriginGroupOverride{
		OriginGroup: &batchRules.ResourceReference{Id: pointer.To(id)},
	}
}
