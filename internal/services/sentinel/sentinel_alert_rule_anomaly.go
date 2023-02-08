package sentinel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type AnomalyRuleRequiredDataConnectorModel struct {
	ConnectorId string   `tfschema:"connector_id"`
	DataTypes   []string `tfschema:"data_types"`
}

// the service always return a fixed one no matter what id we pass, tracked on https://github.com/Azure/azure-rest-api-specs/issues/22485
func AlertRuleAnomalyReadWithPredicate(ctx context.Context, client *securityinsight.SecurityMLAnalyticsSettingsClient, workspaceId workspaces.WorkspaceId, predicateFunc func(v *securityinsight.AnomalySecurityMLAnalyticsSettings) bool) (*securityinsight.AnomalySecurityMLAnalyticsSettings, error) {
	resp, err := client.ListComplete(ctx, workspaceId.ResourceGroupName, workspaceId.WorkspaceName)
	if err != nil {
		return nil, fmt.Errorf("retrieving: %+v", err)
	}

	for resp.NotDone() {
		item := resp.Value()
		if v, ok := item.AsAnomalySecurityMLAnalyticsSettings(); ok {
			if predicateFunc(v) {
				return v, nil
			}

		}
		if err := resp.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("listing next: %+v", err)
		}
	}
	return nil, nil
}

// when the id of workspace is too long, the service return without workspace name:
// "/subscriptions/{sub_id}/resourceGroups/{rg_name}/providers/Microsoft.OperationalInsights/workspaces//providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings/5020e404-9768-4364-98f6-679940c21362",
// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22500
func AlertRuleAnomalyIdFromWorkspaceId(workspaceId workspaces.WorkspaceId, name string) string {
	return parse.NewMLAnalyticsSettingsID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name).ID()
}

func flattenSentinelAlertRuleAnomalyCustomizableObservations(input interface{}) (string, error) {
	value := ""
	val, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to json: %+v", err)
	}
	value = string(val)

	return value, nil
}
