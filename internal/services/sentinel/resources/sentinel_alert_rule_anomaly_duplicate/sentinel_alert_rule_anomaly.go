// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_duplicate

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
)

// when the id of workspace is too long, the service return without workspace name:
// "/subscriptions/{sub_id}/resourceGroups/{rg_name}/providers/Microsoft.OperationalInsights/workspaces//providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings/5020e404-9768-4364-98f6-679940c21362",
// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22500
func AlertRuleAnomalyIdFromWorkspaceId(workspaceId workspaces.WorkspaceId, name string) string {
	return securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name).ID()
}
