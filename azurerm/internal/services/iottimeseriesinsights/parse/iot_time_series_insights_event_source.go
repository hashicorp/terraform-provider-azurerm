package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TimeSeriesInsightsEventSourceId struct {
	ResourceGroup   string
	Name            string
	EnvironmentName string
}

func TimeSeriesInsightsEventSourceID(input string) (*TimeSeriesInsightsEventSourceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Event Source ID %q: %+v", input, err)
	}

	service := TimeSeriesInsightsEventSourceId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.EnvironmentName, err = id.PopSegment("environments"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("eventsources"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
