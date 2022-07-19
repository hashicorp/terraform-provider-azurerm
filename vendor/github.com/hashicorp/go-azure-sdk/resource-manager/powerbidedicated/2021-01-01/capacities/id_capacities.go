package capacities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CapacitiesId{}

// CapacitiesId is a struct representing the Resource ID for a Capacities
type CapacitiesId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DedicatedCapacityName string
}

// NewCapacitiesID returns a new CapacitiesId struct
func NewCapacitiesID(subscriptionId string, resourceGroupName string, dedicatedCapacityName string) CapacitiesId {
	return CapacitiesId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DedicatedCapacityName: dedicatedCapacityName,
	}
}

// ParseCapacitiesID parses 'input' into a CapacitiesId
func ParseCapacitiesID(input string) (*CapacitiesId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacitiesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacitiesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DedicatedCapacityName, ok = parsed.Parsed["dedicatedCapacityName"]; !ok {
		return nil, fmt.Errorf("the segment 'dedicatedCapacityName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCapacitiesIDInsensitively parses 'input' case-insensitively into a CapacitiesId
// note: this method should only be used for API response data and not user input
func ParseCapacitiesIDInsensitively(input string) (*CapacitiesId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacitiesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacitiesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DedicatedCapacityName, ok = parsed.Parsed["dedicatedCapacityName"]; !ok {
		return nil, fmt.Errorf("the segment 'dedicatedCapacityName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCapacitiesID checks that 'input' can be parsed as a Capacities ID
func ValidateCapacitiesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapacitiesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capacities ID
func (id CapacitiesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.PowerBIDedicated/capacities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DedicatedCapacityName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capacities ID
func (id CapacitiesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPowerBIDedicated", "Microsoft.PowerBIDedicated", "Microsoft.PowerBIDedicated"),
		resourceids.StaticSegment("staticCapacities", "capacities", "capacities"),
		resourceids.UserSpecifiedSegment("dedicatedCapacityName", "dedicatedCapacityValue"),
	}
}

// String returns a human-readable description of this Capacities ID
func (id CapacitiesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dedicated Capacity Name: %q", id.DedicatedCapacityName),
	}
	return fmt.Sprintf("Capacities (%s)", strings.Join(components, "\n"))
}
