package frontdoors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
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

func (id FrontDoorId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door", segmentsStr)
}

func (id FrontDoorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ParseFrontDoorID parses a FrontDoor ID into an FrontDoorId struct
func ParseFrontDoorID(input string) (*FrontDoorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("frontDoors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseFrontDoorIDInsensitively parses an FrontDoor ID into an FrontDoorId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseFrontDoorID method should be used instead for validation etc.
func ParseFrontDoorIDInsensitively(input string) (*FrontDoorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontDoors' segment
	frontDoorsKey := "frontDoors"
	for key := range id.Path {
		if strings.EqualFold(key, frontDoorsKey) {
			frontDoorsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(frontDoorsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
