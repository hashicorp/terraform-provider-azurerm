package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FrontDoorId struct {
	SubscriptionId string
	Name           string
	ResourceGroup  string
}

func NewFrontDoorID(subscriptionId, resourceGroup, name string) FrontDoorId {
	return FrontDoorId{
		Name:          name,
		ResourceGroup: resourceGroup,
	}
}

func FrontDoorIDInsensitively(input string) (*FrontDoorId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor ID %q: %+v", input, err)
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return frontDoorId, nil
}

func FrontDoorID(input string) (*FrontDoorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor ID %q: %+v", input, err)
	}

	frontDoorId := FrontDoorId{
		ResourceGroup: id.ResourceGroup,
	}

	if frontDoorId.Name, err = id.PopSegment("frontDoors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &frontDoorId, nil
}

func (id FrontDoorId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func parseFrontDoorChildResourceId(input string) (*FrontDoorId, *azure.ResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, nil, err
	}

	frontdoor := FrontDoorId{
		ResourceGroup: id.ResourceGroup,
	}

	for key, value := range id.Path {
		// In Azure API's should follow Postel's Law - where URI's should be insensitive for requests,
		// but case-sensitive when referencing URI's in responses. Unfortunately the Networking API's
		// treat both as case-insensitive - so until these API's follow the spec we need to identify
		// the correct casing here.
		if strings.EqualFold(key, "frontDoors") {
			frontdoor.Name = value
			delete(id.Path, key)
			break
		}
	}

	return &frontdoor, id, nil
}
