package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TimeSeriesInsightsAccessPolicyId struct {
	ResourceGroup   string
	Name            string
	EnvironmentName string
}

func TimeSeriesInsightsAccessPolicyID(input string) (*TimeSeriesInsightsAccessPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Access Policy ID %q: %+v", input, err)
	}

	service := TimeSeriesInsightsAccessPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.EnvironmentName, err = id.PopSegment("environments"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("accesspolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
