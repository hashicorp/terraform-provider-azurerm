package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsStorageInsightsId struct {
	SubscriptionId           string
	ResourceGroup            string
	WorkspaceName            string
	StorageInsightConfigName string
}

func NewLogAnalyticsStorageInsightsID(subscriptionId, resourceGroup, workspaceName, storageInsightConfigName string) LogAnalyticsStorageInsightsId {
	return LogAnalyticsStorageInsightsId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		WorkspaceName:            workspaceName,
		StorageInsightConfigName: storageInsightConfigName,
	}
}

func (id LogAnalyticsStorageInsightsId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Storage Insight Config Name %q", id.StorageInsightConfigName),
	}
	return strings.Join(segments, " / ")
}

func (id LogAnalyticsStorageInsightsId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/storageInsightConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.StorageInsightConfigName)
}

// LogAnalyticsStorageInsightsID parses a LogAnalyticsStorageInsights ID into an LogAnalyticsStorageInsightsId struct
func LogAnalyticsStorageInsightsID(input string) (*LogAnalyticsStorageInsightsId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsStorageInsightsId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.StorageInsightConfigName, err = id.PopSegment("storageInsightConfigs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
