package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AvailabilitySetId struct {
	ResourceGroup string
	Name          string
}

func AvailabilitySetID(input string) (*AvailabilitySetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Availability Set ID %q: %+v", input, err)
	}

	set := AvailabilitySetId{
		ResourceGroup: id.ResourceGroup,
	}

	if set.Name, err = id.PopSegment("availabilitySets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &set, nil
}
