package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsDataSourceWindowsEventId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
	DataSourceName string
}

func NewLogAnalyticsDataSourceWindowsEventID(subscriptionId, resourceGroup, workspaceName, dataSourceName string) LogAnalyticsDataSourceWindowsEventId {
	return LogAnalyticsDataSourceWindowsEventId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
		DataSourceName: dataSourceName,
	}
}

func (id LogAnalyticsDataSourceWindowsEventId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Data Source Name %q", id.DataSourceName),
	}
	return strings.Join(segments, " / ")
}

func (id LogAnalyticsDataSourceWindowsEventId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/dataSources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.DataSourceName)
}

// LogAnalyticsDataSourceWindowsEventID parses a LogAnalyticsDataSourceWindowsEvent ID into an LogAnalyticsDataSourceWindowsEventId struct
func LogAnalyticsDataSourceWindowsEventID(input string) (*LogAnalyticsDataSourceWindowsEventId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsDataSourceWindowsEventId{
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
	if resourceId.DataSourceName, err = id.PopSegment("dataSources"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
