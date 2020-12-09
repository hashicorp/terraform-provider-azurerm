package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsDataExportId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
	DataexportName string
}

func NewLogAnalyticsDataExportID(subscriptionId, resourceGroup, workspaceName, dataexportName string) LogAnalyticsDataExportId {
	return LogAnalyticsDataExportId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
		DataexportName: dataexportName,
	}
}

func (id LogAnalyticsDataExportId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Dataexport Name %q", id.DataexportName),
	}
	return strings.Join(segments, " / ")
}

func (id LogAnalyticsDataExportId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/dataexports/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.DataexportName)
}

// LogAnalyticsDataExportID parses a LogAnalyticsDataExport ID into an LogAnalyticsDataExportId struct
func LogAnalyticsDataExportID(input string) (*LogAnalyticsDataExportId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsDataExportId{
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
	if resourceId.DataexportName, err = id.PopSegment("dataexports"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
