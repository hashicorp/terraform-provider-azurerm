// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
)

// TestResolveLogAnalyticsWorkspaceId covers the state-resolution logic that fixes
// https://github.com/hashicorp/terraform-provider-azurerm/issues/32705, where a
// `log_analytics_workspace_id` pointing at a Workspace in a *different* subscription to the
// Container App Environment produced a perpetual diff. The API only ever returns a `customerId`
// GUID, so the provider has to resolve that back to an ARM resource ID by listing Workspaces in
// its own configured subscription (`findWorkspaceResourceIDFromCustomerID`) - which can never
// succeed for a cross-subscription Workspace. In that case the previously known value must be
// preserved instead of being overwritten with an empty string.
func TestResolveLogAnalyticsWorkspaceId(t *testing.T) {
	crossSubWorkspaceId := "/subscriptions/aaaaaaaa-0000-0000-0000-aaaaaaaaaaaa/resourceGroups/rg-myapp-dev/providers/Microsoft.App/managedEnvironments/cenv-myapp-dev"
	resolvedId := workspaces.NewWorkspaceID("bbbbbbbb-0000-0000-0000-bbbbbbbbbbbb", "rg-monitoring-shared", "log-central-mgt")

	testCases := []struct {
		name          string
		workspaceId   *workspaces.WorkspaceId
		lookupErr     error
		existingValue string
		expected      string
	}{
		{
			name:          "workspace resolved successfully in the current subscription",
			workspaceId:   &resolvedId,
			lookupErr:     nil,
			existingValue: "",
			expected:      resolvedId.ID(),
		},
		{
			name:          "workspace resolved takes precedence over existing state",
			workspaceId:   &resolvedId,
			lookupErr:     nil,
			existingValue: crossSubWorkspaceId,
			expected:      resolvedId.ID(),
		},
		{
			name:          "no matching workspace found in current subscription: fall back to existing state",
			workspaceId:   nil,
			lookupErr:     nil,
			existingValue: crossSubWorkspaceId,
			expected:      crossSubWorkspaceId,
		},
		{
			// This is the exact scenario from #32705: the provider's configured subscription has
			// no Log Analytics Workspaces of its own at all (they all live in a shared/monitoring
			// subscription), so `findWorkspaceResourceIDFromCustomerID` returns an error rather
			// than `nil, nil`. Before this fix that error caused the state value to be silently
			// dropped (persisted as ""), which is what produced the perpetual diff.
			name:          "lookup errors because current subscription has no workspaces at all: fall back to existing state",
			workspaceId:   nil,
			lookupErr:     errors.New("could not resolve Log Analytics Workspace ID for customer-id, no Log Analytics Workspaces found in subscription"),
			existingValue: crossSubWorkspaceId,
			expected:      crossSubWorkspaceId,
		},
		{
			name:          "lookup errors and there is no existing state: value stays empty",
			workspaceId:   nil,
			lookupErr:     errors.New("some transient API error"),
			existingValue: "",
			expected:      "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := resolveLogAnalyticsWorkspaceId(tc.workspaceId, tc.lookupErr, tc.existingValue)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}
