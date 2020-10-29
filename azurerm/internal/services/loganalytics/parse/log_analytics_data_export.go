package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsDataExportId struct {
	ResourceGroup string
	WorkspaceName string
	WorkspaceID   string
	Name          string
}

func LogAnalyticsDataExportID(input string) (*LogAnalyticsDataExportId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Log Analytics Data Export Rule ID %q: %+v", input, err)
	}

	logAnalyticsDataExport := LogAnalyticsDataExportId{
		ResourceGroup: id.ResourceGroup,
	}

	if logAnalyticsDataExport.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if logAnalyticsDataExport.WorkspaceID = fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/%s/workspaces/%s", id.SubscriptionID, id.ResourceGroup, id.Provider, logAnalyticsDataExport.WorkspaceName); err != nil {
		return nil, fmt.Errorf("formatting Log Analytics Data Export Rule workspace ID %q", input)
	}
	if logAnalyticsDataExport.Name, err = id.PopSegment("dataExports"); err != nil {
		// API Issue the casing changes for the ID
		if logAnalyticsDataExport.Name, err = id.PopSegment("dataexports"); err != nil {
			return nil, err
		}
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logAnalyticsDataExport, nil
}
