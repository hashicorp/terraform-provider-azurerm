package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProximityPlacementGroupId struct {
	ResourceGroup string
	Name          string
}

func NewProximityPlacementGroupId(resourceGroup, name string) ProximityPlacementGroupId {
	return ProximityPlacementGroupId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id ProximityPlacementGroupId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/proximityPlacementGroups/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func ProximityPlacementGroupID(input string) (*ProximityPlacementGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Proximity Placement Group ID %q: %+v", input, err)
	}

	server := ProximityPlacementGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("proximityPlacementGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
