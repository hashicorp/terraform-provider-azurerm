package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TimeSeriesInsightsEnvironmentId struct {
	ResourceGroup string
	Name          string
}

func TimeSeriesInsightsEnvironmentID(input string) (*TimeSeriesInsightsEnvironmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Environment ID %q: %+v", input, err)
	}

	service := TimeSeriesInsightsEnvironmentId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("environments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
