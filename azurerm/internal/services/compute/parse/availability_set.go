package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AvailabilitySetId struct {
	ResourceGroup string
	Name          string
}

func NewAvailabilitySetId(resourceGroup, name string) AvailabilitySetId {
	return AvailabilitySetId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id AvailabilitySetId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/availabilitySets/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func AvailabilitySetID(input string) (*AvailabilitySetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Availability Set ID %q: %+v", input, err)
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
