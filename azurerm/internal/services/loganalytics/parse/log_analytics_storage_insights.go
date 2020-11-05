package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsStorageInsightsId struct {
	ResourceGroup string
	WorkspaceName string
	WorkspaceID   string
	Name          string
}

func LogAnalyticsStorageInsightsID(input string) (*LogAnalyticsStorageInsightsId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Log Analytics Storage Insights ID %q: %+v", input, err)
	}

	logAnalyticsStorageInsight := LogAnalyticsStorageInsightsId{
		ResourceGroup: id.ResourceGroup,
	}
	if logAnalyticsStorageInsight.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if logAnalyticsStorageInsight.WorkspaceID = fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/%s/workspaces/%s", id.SubscriptionID, id.ResourceGroup, id.Provider, logAnalyticsStorageInsight.WorkspaceName); err != nil {
		return nil, fmt.Errorf("formatting Log Analytics Storage Insights workspace ID %q", input)
	}
	if logAnalyticsStorageInsight.Name, err = id.PopSegment("storageInsightConfigs"); err != nil {
		if logAnalyticsStorageInsight.Name, err = id.PopSegment("storageinsightconfigs"); err != nil {
			return nil, err
		}
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logAnalyticsStorageInsight, nil
}
