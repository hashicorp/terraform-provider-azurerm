package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsStorageInsightConfigId struct {
	ResourceGroup string
	WorkspaceName string
	WorkspaceID   string
	Name          string
}

func LogAnalyticsStorageInsightConfigID(input string) (*LogAnalyticsStorageInsightConfigId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Log Analytics Storage Insight Config ID %q: %+v", input, err)
	}

	logAnalyticsStorageInsightConfig := LogAnalyticsStorageInsightConfigId{
		ResourceGroup: id.ResourceGroup,
	}
	if logAnalyticsStorageInsightConfig.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if logAnalyticsStorageInsightConfig.WorkspaceID = fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/%s/workspaces/%s", id.SubscriptionID, id.ResourceGroup, id.Provider, logAnalyticsStorageInsightConfig.WorkspaceName); err != nil {
		return nil, fmt.Errorf("formatting Log Analytics Storage Insight Config workspace ID %q", input)
	}
	if logAnalyticsStorageInsightConfig.Name, err = id.PopSegment("storageInsightConfigs"); err != nil {
		if logAnalyticsStorageInsightConfig.Name, err = id.PopSegment("storageinsightconfigs"); err != nil {
			return nil, err
		}
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logAnalyticsStorageInsightConfig, nil
}
