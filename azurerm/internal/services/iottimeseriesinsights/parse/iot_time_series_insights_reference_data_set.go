package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TimeSeriesInsightsReferenceDataSetId struct {
	ResourceGroup   string
	EnvironmentName string
	Name            string
}

func TimeSeriesInsightsReferenceDataSetID(input string) (*TimeSeriesInsightsReferenceDataSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Reference Dataset ID %q: %+v", input, err)
	}

	service := TimeSeriesInsightsReferenceDataSetId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.EnvironmentName, err = id.PopSegment("environments"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("referenceDataSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
