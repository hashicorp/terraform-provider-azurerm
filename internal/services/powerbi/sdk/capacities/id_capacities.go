package capacities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CapacitiesId struct {
	SubscriptionId string
	ResourceGroup  string
	CapacityName   string
}

func NewCapacitiesID(subscriptionId, resourceGroup, capacityName string) CapacitiesId {
	return CapacitiesId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		CapacityName:   capacityName,
	}
}

func (id CapacitiesId) String() string {
	segments := []string{
		fmt.Sprintf("Capacity Name %q", id.CapacityName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Capacities", segmentsStr)
}

func (id CapacitiesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.PowerBIDedicated/capacities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CapacityName)
}

// ParseCapacitiesID parses a Capacities ID into an CapacitiesId struct
func ParseCapacitiesID(input string) (*CapacitiesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CapacitiesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CapacityName, err = id.PopSegment("capacities"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseCapacitiesIDInsensitively parses an Capacities ID into an CapacitiesId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseCapacitiesID method should be used instead for validation etc.
func ParseCapacitiesIDInsensitively(input string) (*CapacitiesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CapacitiesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'capacities' segment
	capacitiesKey := "capacities"
	for key := range id.Path {
		if strings.EqualFold(key, capacitiesKey) {
			capacitiesKey = key
			break
		}
	}
	if resourceId.CapacityName, err = id.PopSegment(capacitiesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
