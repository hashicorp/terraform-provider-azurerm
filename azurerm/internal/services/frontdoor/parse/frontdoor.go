package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// TODO: tests

type FrontDoorId struct {
	Name          string
	ResourceGroup string
}

func NewFrontDoorID(resourceGroup, name string) FrontDoorId {
	return FrontDoorId{
		Name:          name,
		ResourceGroup: resourceGroup,
	}
}

func FrontDoorID(input string) (*FrontDoorId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse FrontDoor ID %q: %+v", input, err)
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return frontDoorId, nil
}

func (id FrontDoorId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontdoors/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func parseFrontDoorChildResourceId(input string) (*FrontDoorId, *azure.ResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, nil, err
	}

	frontdoor := FrontDoorId{
		ResourceGroup: id.ResourceGroup,
	}

	// TODO: ensure this is Normalized, presumably to frontdoor
	// resourceGroup := id.ResourceGroup
	//	name := id.Path["frontdoors"]
	//	// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
	//	if name == "" {
	//		name = id.Path["Frontdoors"]
	//	}

	if frontdoor.Name, err = id.PopSegment("frontdoors"); err != nil {
		return nil, nil, err
	}

	return &frontdoor, id, nil
}
