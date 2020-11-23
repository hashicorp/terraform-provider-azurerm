package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FrontDoorId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewFrontDoorID(subscriptionId, resourceGroup, name string) FrontDoorId {
	return FrontDoorId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func FrontDoorID(input string) (*FrontDoorId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor ID %q: %+v", input, err)
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return frontDoorId, nil
}

func FrontDoorIDForImport(input string) (*FrontDoorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor ID %q: %+v", input, err)
	}

	frontDoorId := FrontDoorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if frontDoorId.Name, err = id.PopSegment("frontDoors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &frontDoorId, nil
}

func (id FrontDoorId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func parseFrontDoorChildResourceId(input string) (*FrontDoorId, *azure.ResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, nil, err
	}

	frontdoor := FrontDoorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
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
