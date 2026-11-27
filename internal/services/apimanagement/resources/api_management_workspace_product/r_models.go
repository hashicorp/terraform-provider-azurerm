// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_product

type ApiManagementWorkspaceProductModel struct {
	Name                       string `tfschema:"name"`
	ApiManagementWorkspaceId   string `tfschema:"api_management_workspace_id"`
	RequireApprovalEnabled     bool   `tfschema:"require_approval_enabled"`
	Description                string `tfschema:"description"`
	DisplayName                string `tfschema:"display_name"`
	PublishedEnabled           bool   `tfschema:"published_enabled"`
	RequireSubscriptionEnabled bool   `tfschema:"require_subscription_enabled"`
	SubscriptionsLimit         int64  `tfschema:"subscriptions_limit"`
	Terms                      string `tfschema:"terms"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceProductModel{}
}
