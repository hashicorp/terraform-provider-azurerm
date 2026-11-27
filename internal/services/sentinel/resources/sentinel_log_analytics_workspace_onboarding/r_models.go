// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_log_analytics_workspace_onboarding

type SecurityInsightsSentinelOnboardingStateModel struct {
	CustomerManagedKeyEnabled bool   `tfschema:"customer_managed_key_enabled"`
	WorkspaceId               string `tfschema:"workspace_id"`
}

func (r Resource) ModelObject() any {
	return &SecurityInsightsSentinelOnboardingStateModel{}
}
