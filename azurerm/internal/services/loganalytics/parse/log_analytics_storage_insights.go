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

// Note - this API currently lower-cases all return values
// Issue tracked here - https://github.com/Azure/azure-sdk-for-go/issues/13268
const fmtWorkspaceId = "/subscriptions/%s/resourcegroups/%s/providers/%s/workspaces/%s"

func (id LogAnalyticsStorageInsightsId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourcegroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/storageInsightConfigs/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.WorkspaceName, id.Name)
}

func NewLogAnalyticsStorageInsightsId(resourceGroup, workspaceId, name string) LogAnalyticsStorageInsightsId {
	// (@jackofallops) ignoring error here as already passed through validation in schema
	workspace, _ := LogAnalyticsWorkspaceID(workspaceId)
	return LogAnalyticsStorageInsightsId{
		ResourceGroup: resourceGroup,
		WorkspaceName: workspace.Name,
		WorkspaceID:   workspaceId,
		Name:          name,
	}
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
	if logAnalyticsStorageInsight.WorkspaceID = fmt.Sprintf(fmtWorkspaceId, id.SubscriptionID, id.ResourceGroup, id.Provider, logAnalyticsStorageInsight.WorkspaceName); err != nil {
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
