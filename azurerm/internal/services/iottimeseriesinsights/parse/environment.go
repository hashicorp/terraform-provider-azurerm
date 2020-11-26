package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EnvironmentId struct {
	ResourceGroup string
	Name          string
}

func NewEnvironmentID(resourceGroup, name string) EnvironmentId {
	return EnvironmentId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id EnvironmentId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.TimeSeriesInsights/environments/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func EnvironmentID(input string) (*EnvironmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Environment ID %q: %+v", input, err)
	}

	service := EnvironmentId{
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
