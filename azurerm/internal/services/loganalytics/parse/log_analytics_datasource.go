package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsDataSourceId struct {
	ResourceGroup string
	Workspace     string
	Name          string
}

func LogAnalyticsDataSourceID(input string) (*LogAnalyticsDataSourceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Log Analytics Data Source ID %q: %+v", input, err)
	}

	server := LogAnalyticsDataSourceId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Workspace, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if server.Name, err = id.PopSegment("datasources"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
