package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ReferenceDataSetId struct {
	ResourceGroup   string
	EnvironmentName string
	Name            string
}

func ReferenceDataSetID(input string) (*ReferenceDataSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Reference Dataset ID %q: %+v", input, err)
	}

	service := ReferenceDataSetId{
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
